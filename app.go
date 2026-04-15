package compogo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Compogo/compogo/closer"
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/configurator"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/compogo/logger"
	hashSlice "github.com/Compogo/types/hash_slice"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/set"
)

// App represents the main application container.
// It manages component lifecycle, dependencies, and graceful shutdown.
type App struct {
	name string

	configCmp *component.Component
	config    *Config

	containerCmp *component.Component
	container    container.Container

	configuratorCmp *component.Component
	configurator    configurator.Configurator

	closerCmp *component.Component
	closer    closer.Closer

	loggerCmp *component.Component
	logger    logger.Logger

	components hashSlice.HashSlice[*component.Component]
	wg         sync.WaitGroup
	waitMutex  sync.Mutex

	steps    *linker.Linker[component.Step, set.Set[*component.Component]]
	timeouts *linker.Linker[component.Step, time.Duration]

	isRunning atomic.Bool

	parent *App
}

// NewApp creates a new application instance with the given name and options.
// The config component is automatically added to ensure basic configuration is always present.
func NewApp(name string, options ...Option) *App {
	app := &App{
		name: name,
		steps: linker.NewLinker[component.Step, set.Set[*component.Component]](
			linker.NewLink(component.Init, set.NewSet[*component.Component]()),
			linker.NewLink(component.BindFlag, set.NewSet[*component.Component]()),
			linker.NewLink(component.Configuration, set.NewSet[*component.Component]()),
			linker.NewLink(component.PreExecute, set.NewSet[*component.Component]()),
			linker.NewLink(component.Execute, set.NewSet[*component.Component]()),
			linker.NewLink(component.PostExecute, set.NewSet[*component.Component]()),
			linker.NewLink(component.PreWait, set.NewSet[*component.Component]()),
			linker.NewLink(component.Wait, set.NewSet[*component.Component]()),
			linker.NewLink(component.PostWait, set.NewSet[*component.Component]()),
			linker.NewLink(component.PreStop, set.NewSet[*component.Component]()),
			linker.NewLink(component.Stop, set.NewSet[*component.Component]()),
			linker.NewLink(component.PostStop, set.NewSet[*component.Component]()),
		),
		timeouts: linker.NewLinker[component.Step, time.Duration](),
	}

	options = append(options, withConfig(NewConfig()))

	for _, option := range options {
		option(app)
	}

	return app
}

// AddComponents registers one or more components and their dependencies in the application.
// Components are initialized immediately and cannot be added while the app is running.
// Returns an error if validation fails or the app is already running.
func (app *App) AddComponents(components ...*component.Component) (err error) {
	if err = app.validate(); err != nil {
		return err
	}

	if app.IsRunning() {
		return fmt.Errorf("[compogo][%s] %w", app.name, AppIsRunningError)
	}

	for _, cmp := range components {
		if len(cmp.Dependencies) > 0 {
			if err = app.AddComponents(cmp.Dependencies...); err != nil {
				return err
			}
		}

		if !app.existComponent(cmp) {
			_, _ = app.components.Add(cmp)
		}
	}

	return nil
}

func (app *App) existComponent(component *component.Component) bool {
	if app.parent != nil {
		if exist := app.parent.existComponent(component); exist {
			return true
		}
	}

	return app.components.Contains(component)
}

// BindFlags binds command-line flags for all registered components.
// Must be called before Serve() and cannot be called while the app is running.
func (app *App) BindFlags(flagSet flag.FlagSet) (err error) {
	if err = app.validate(); err != nil {
		return err
	}

	if app.IsRunning() {
		return fmt.Errorf("[compogo][%s].BindFlags: %w", app.name, AppIsRunningError)
	}

	if err = app.runComponents(component.Init); err != nil {
		return fmt.Errorf("[compogo][%s].BindFlags: %w", app.name, err)
	}

	if app.parent != nil {
		if err = app.parent.BindFlags(flagSet); err != nil {
			return err
		}
	}

	components := app.getCoreComponents()
	components = append(components, app.components.ToSlice()...)

	bindFlags, _ := app.steps.Get(component.BindFlag)

	for _, cmp := range components {
		if !bindFlags.Contains(cmp) && cmp.BindFlags != nil {
			if err = cmp.BindFlags(flagSet, app.container); err != nil {
				return err
			}

			bindFlags.Add(cmp)
		}
	}

	return nil
}

// Serve starts the application and runs all components through their lifecycle:
// 1. Sequential execution of PreExecute, Execute, PostExecute, PreWait
// 2. Concurrent execution of all Wait components
// 3. Sequential execution of PostWait
// 4. Wait for shutdown signal or first error
// 5. Sequential execution of PreStop, Stop, PostStop
// 6. Wait for all goroutines to finish
func (app *App) Serve() (err error) {
	if err = app.validate(); err != nil {
		return err
	}

	if app.IsRunning() {
		return fmt.Errorf("[compogo][%s] %w", app.name, AppIsRunningError)
	}

	app.logger.Infof("[compogo][%s] Running", app.name)

	app.setRunning(true)
	defer app.setRunning(false)

	if err = app.runComponents(component.Init, component.Configuration, component.PreExecute, component.Execute, component.PostExecute, component.PreWait); err != nil {
		return err
	}

	components := app.getAllComponents()

	chainErr := make(chan error, 1)
	defer close(chainErr)

	ctx, cancelFunc := context.WithCancel(app.closer.GetContext())
	defer cancelFunc()

	waitComponents, _ := app.steps.Get(component.Wait)

	for _, cmp := range components {
		if cmp.Wait != nil && !waitComponents.Contains(cmp) {
			app.serveComponent(ctx, cmp, waitComponents, chainErr)
			waitComponents.Add(cmp)
		}
	}

	if err = app.runComponents(component.PostWait); err != nil {
		return err
	}

	select {
	case waitErr := <-chainErr:
		err = errors.Join(err, fmt.Errorf("[compogo][%s] %w", app.name, waitErr))

		if closerErr := app.closer.Close(); closerErr != nil {
			err = errors.Join(err, fmt.Errorf("[compogo][%s] %w", app.name, closerErr))
		}
	case <-ctx.Done():
		break
	}

	app.logger.Info("Shutdown")

	if err = app.runComponents(component.PreStop, component.Stop, component.PostStop); err != nil {
		return err
	}

	app.wg.Wait()

	return err
}

func (app *App) runComponents(steps ...component.Step) (err error) {
	if err = app.validate(); err != nil {
		return err
	}

	if app.parent != nil {
		if err = app.parent.runComponents(steps...); err != nil {
			return err
		}
	}

	components := app.getCoreComponents()
	components = append(components, app.components.ToSlice()...)

	for _, step := range steps {
		if err = app.runStepComponents(step, components...); err != nil {
			return err
		}
	}

	return nil
}

func (app *App) runStepComponents(step component.Step, components ...*component.Component) (err error) {
	var ctx context.Context
	var cancelFunc context.CancelFunc
	var fnc component.StepFunc
	timeout := app.timeouts.GetOrDefault(step, time.Second)
	stepComponents, _ := app.steps.Get(step)

	cxtBackground := context.Background()

	errChan := make(chan error, 1)
	defer close(errChan)

	for _, cmp := range components {
		if len(cmp.Dependencies) > 0 {
			if err = app.runStepComponents(step, cmp.Dependencies...); err != nil {
				return err
			}
		}

		if stepComponents.Contains(cmp) {
			continue
		}

		fnc, err = cmp.GetStepFunc(step)
		if err != nil {
			return fmt.Errorf("[compogo][%s][%s]: %w", app.name, step, err)
		}

		if fnc == nil {
			continue
		}

		cmpStepTimeout := timeout
		if step == component.Execute && cmp.ExecuteDuration != nil {
			cmpStepTimeout = cmp.ExecuteDuration()
		}

		if cmpStepTimeout > 0 {
			ctx, cancelFunc = context.WithTimeout(cxtBackground, cmpStepTimeout)
		} else {
			ctx, cancelFunc = context.WithCancel(cxtBackground)
		}

		go func() {
			defer cancelFunc()
			if err := fnc(app.container); err != nil {
				errChan <- err
			}
		}()

		select {
		case <-ctx.Done():
			break
		case err := <-errChan:
			return err
		}

		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return fmt.Errorf("[compogo][%s] component - '%s', step - '%s', failed: %w", app.name, cmp.Name, step, ComponentStepTimeoutError)
		}

		stepComponents.Add(cmp)
	}

	return nil
}

func (app *App) serveComponent(ctx context.Context, cmp *component.Component, waitComponents set.Set[*component.Component], chainErr chan error) {
	app.waitMutex.Lock()
	defer app.waitMutex.Unlock()

	app.wg.Add(1)
	go func(ctx context.Context, cmp *component.Component, waitComponents set.Set[*component.Component]) {
		defer app.wg.Done()
		defer func() {
			app.waitMutex.Lock()
			defer app.waitMutex.Unlock()
			waitComponents.Remove(cmp)
		}()

		defer func() {
			if r := recover(); r != nil {
				chainErr <- fmt.Errorf("[compogo][%s] component '%s' wait failed: %s", app.name, cmp.Name, r)
			}
		}()

		ctx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()

		if err := cmp.Wait(ctx, app.container); err != nil {
			chainErr <- fmt.Errorf("[compogo][%s] component '%s' wait failed: %w", app.name, cmp.Name, err)
		}
	}(ctx, cmp, waitComponents)
}

func (app *App) getAllComponents() []*component.Component {
	components := app.getCoreComponents()

	if app.parent != nil {
		components = append(components, app.parent.getAllComponents()...)
	}

	components = append(components, app.components.ToSlice()...)

	return components
}

func (app *App) setRunning(isRunning bool) {
	if app.parent != nil {
		app.parent.setRunning(isRunning)
	}

	app.isRunning.Store(isRunning)
}

// IsRunning returns true if the application is in the running state.
func (app *App) IsRunning() bool {
	if app.parent != nil {
		if isRunning := app.parent.IsRunning(); isRunning {
			return true
		}
	}

	return app.isRunning.Load()
}

func (app *App) validate() error {
	var err error

	if app.parent != nil {
		if err = app.parent.validate(); err != nil {
			return err
		}
	}

	if app.container == nil {
		err = fmt.Errorf("%w%w\n", err, ContainerUndefinedError)
	}

	if app.configurator == nil {
		err = fmt.Errorf("%w%w\n", err, ConfiguratorUndefinedError)
	}

	if app.closer == nil {
		err = fmt.Errorf("%w%w\n", err, CloserUndefinedError)
	}

	if app.logger == nil {
		err = fmt.Errorf("%w%w\n", err, LoggerUndefinedError)
	}

	if err != nil {
		err = fmt.Errorf("[compogo][%s]%w", app.name, err)
	}

	return err
}

// Clone creates a child application that inherits all core services
// (config, container, configurator, closer) but has its own logger and component set.
// Useful for creating isolated sub-applications (e.g., for testing or modules).
func (app *App) Clone(name string) *App {
	return &App{
		name:         fmt.Sprintf("%s.%s", app.name, name),
		config:       app.config,
		container:    app.container,
		configurator: app.configurator,
		closer:       app.closer,
		logger:       app.logger.GetLogger(name),
		parent:       app,
	}
}

func (app *App) getCoreComponents() []*component.Component {
	if app.configCmp == nil {
		return nil
	}

	return []*component.Component{
		app.configuratorCmp,
		app.containerCmp,
		app.configCmp,
		app.closerCmp,
		app.loggerCmp,
	}
}
