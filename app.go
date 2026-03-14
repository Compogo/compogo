package compogo

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Compogo/compogo/closer"
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/configurator"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/compogo/logger"
	"github.com/Compogo/compogo/types"
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

	components types.Set[*component.Component]
	wg         sync.WaitGroup

	bindFlags types.Set[*component.Component]

	init          types.Set[*component.Component]
	configuration types.Set[*component.Component]

	preRun  types.Set[*component.Component]
	run     types.Set[*component.Component]
	postRun types.Set[*component.Component]

	preWait   types.Set[*component.Component]
	wait      types.Set[*component.Component]
	waitMutex sync.Mutex
	postWait  types.Set[*component.Component]

	preStop  types.Set[*component.Component]
	stop     types.Set[*component.Component]
	postStop types.Set[*component.Component]

	isRunning atomic.Bool

	parent *App
}

// NewApp creates a new application instance with the given name and options.
// The config component is automatically added to ensure basic configuration is always present.
func NewApp(name string, options ...Option) *App {
	app := &App{name: name}

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
			app.components.Add(cmp)
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

	for _, cmp := range components {
		if !app.bindFlags.Contains(cmp) && cmp.BindFlags != nil {
			if err = cmp.BindFlags(flagSet, app.container); err != nil {
				return err
			}

			app.bindFlags.Add(cmp)
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

	for _, cmp := range components {
		app.serveComponent(ctx, cmp, chainErr)
	}

	if err = app.runComponents(component.PostWait); err != nil {
		return err
	}

	select {
	case waitErr := <-chainErr:
		err = fmt.Errorf("%w\n[compogo][%s] %w", err, app.name, waitErr)

		if closerErr := app.closer.Close(); closerErr != nil {
			err = fmt.Errorf("%w\n[compogo][%s] %w", err, app.name, closerErr)
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

	errChan := make(chan error, 1)
	defer close(errChan)

	for _, cmp := range components {
		if len(cmp.Dependencies) > 0 {
			if err = app.runStepComponents(step, cmp.Dependencies...); err != nil {
				return err
			}
		}

		switch step {
		// init
		case component.Init:
			if app.init.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.InitDuration)
			fnc = cmp.Init

			app.init.Add(cmp)

		case component.Configuration:
			if app.configuration.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.ConfigurationDuration)
			fnc = cmp.Configuration

			app.configuration.Add(cmp)

		// run
		case component.PreExecute:
			if app.preRun.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.PreRunDuration)
			fnc = cmp.PreExecute

			app.preRun.Add(cmp)
		case component.Execute:
			if app.run.Contains(cmp) {
				continue
			}

			duration := app.config.RunDuration
			if cmp.ExecuteDuration != nil {
				duration = cmp.ExecuteDuration()
			}

			if duration > 0 {
				ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), duration)
			} else {
				ctx, cancelFunc = context.WithCancel(app.closer.GetContext())
			}

			fnc = cmp.Execute

			app.run.Add(cmp)
		case component.PostExecute:
			if app.postRun.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.PostRunDuration)
			fnc = cmp.PostExecute

			app.postRun.Add(cmp)
		// wait
		case component.PreWait:
			if app.preWait.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.PreWaitDuration)
			fnc = cmp.PreWait

			app.preWait.Add(cmp)
		case component.PostWait:
			if app.postWait.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.PostWaitDuration)
			fnc = cmp.PostWait

			app.postWait.Add(cmp)
		// stop
		case component.PreStop:
			if app.preStop.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.PreStopDuration)
			fnc = cmp.PreStop

			app.preStop.Add(cmp)
		case component.Stop:
			if app.stop.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.StopDuration)
			fnc = cmp.Stop

			app.stop.Add(cmp)
		case component.PostStop:
			if app.postStop.Contains(cmp) {
				continue
			}

			ctx, cancelFunc = context.WithTimeout(app.closer.GetContext(), app.config.PostStopDuration)
			fnc = cmp.PostStop

			app.postStop.Add(cmp)
		default:
			return fmt.Errorf("[compogo][%s][%s]: %w", app.name, step, StepUndefinedError)
		}

		if fnc == nil {
			cancelFunc()
			continue
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
	}

	return nil
}

func (app *App) serveComponent(ctx context.Context, cmp *component.Component, chainErr chan error) {
	app.waitMutex.Lock()
	defer app.waitMutex.Unlock()

	if cmp.Wait == nil || app.wait.Contains(cmp) {
		return
	}

	app.wg.Add(1)
	app.wait.Add(cmp)
	go func(ctx context.Context, cmp *component.Component) {
		defer app.wg.Done()
		defer func() {
			app.waitMutex.Lock()
			defer app.waitMutex.Unlock()
			app.wait.Remove(cmp)
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
	}(ctx, cmp)
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
