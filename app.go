package compogo

import (
	"fmt"
	"os"
	"slices"
	"sync/atomic"
	"time"

	hashSlice "github.com/Compogo/types/hash_slice"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/set"
	"golang.org/x/sync/errgroup"
)

type App struct {
	name   string
	config *Config

	containerCmp *Component
	container    Container

	configuratorCmp *Component
	configurator    Configurator

	closerCmp *Component
	closer    Closer

	loggerCmp *Component
	logger    Logger

	flagCmp *Component

	components *hashSlice.HashSlice[*Component]
	steps      *linker.Linker[Step, set.Set[*Component]]
	wgErr      errgroup.Group

	parent   *App
	children []*App

	isRunning atomic.Bool
}

func NewApp(name string, options ...Option) (*App, error) {
	app := &App{
		name: name,
		config: &Config{
			Name: name,
			PID:  uint64(os.Getpid()),
		},
		steps: linker.NewLinker[Step, set.Set[*Component]](
			linker.Link(Init, set.NewSet[*Component]()),
			linker.Link(BindFlag, set.NewSet[*Component]()),
			linker.Link(Configuration, set.NewSet[*Component]()),
			linker.Link(PreExecute, set.NewSet[*Component]()),
			linker.Link(Execute, set.NewSet[*Component]()),
			linker.Link(PostExecute, set.NewSet[*Component]()),
			linker.Link(PreWait, set.NewSet[*Component]()),
			linker.Link(Wait, set.NewSet[*Component]()),
			linker.Link(PostWait, set.NewSet[*Component]()),
			linker.Link(PreStop, set.NewSet[*Component]()),
			linker.Link(Stop, set.NewSet[*Component]()),
			linker.Link(PostStop, set.NewSet[*Component]()),
		),
		components: hashSlice.NewHashSlice[*Component](),
		flagCmp:    &flagComponent,
	}

	for _, option := range options {
		option(app)
	}

	if app.container == nil {
		return nil, ContainerUndefinedErr
	}

	if app.configurator == nil {
		return nil, ConfiguratorUndefinedErr
	}

	if app.closer == nil {
		return nil, CloserUndefinedErr
	}

	if app.logger == nil {
		return nil, LoggerUndefinedErr
	}

	if err := app.container.Provide(func() *Config { return app.config }); err != nil {
		return nil, err
	}

	app.logger = app.logger.GetLogger("Compogo").GetLogger(app.name)

	return app, nil
}

func (app *App) AddComponents(components ...*Component) (err error) {
	if app.IsRunning() {
		return fmt.Errorf("[compogo][%s] AddComponents: %w", app.name, AppIsRunningErr)
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

func (app *App) existComponent(component *Component) bool {
	if app.parent != nil {
		if exist := app.parent.existComponent(component); exist {
			return true
		}
	}

	return app.components.Contains(component)
}

// BindFlags выполняет привязку флагов командной строки для всех компонентов.
// Метод рекурсивно обходит все компоненты и зависимости, вызывая их BindFlags-функции.
//
// Порядок выполнения:
//   - Сначала выполняются Init-функции компонентов (для регистрации в DI-контейнере)
//   - Затем в компонентах с учётом зависимостей выполняются BindFlags
//
// Каждый компонент выполняется только один раз (отслеживается через steps).
//
// Пример:
//
//	flagSet := pflag.NewFlagSet("myapp", pflag.ExitOnError)
//	if err := app.BindFlags(flagSet); err != nil {
//	    log.Fatal(err)
//	}
//	flagSet.Parse(os.Args[1:])
func (app *App) BindFlags() error {
	if app.IsRunning() {
		return fmt.Errorf("[compogo][%s] BindFlags: %w", app.name, AppIsRunningErr)
	}

	return app.runSteps(Init, BindFlag)
}

// Serve запускает приложение и все его компоненты.
// Это основной метод, который блокирует выполнение до завершения всех компонентов.
//
// Порядок выполнения этапов:
//   - Init (инициализация компонентов)
//   - Configuration (загрузка конфигурации)
//   - PreExecute → Execute → PostExecute (код однопоточного выполнения)
//   - PreWait → Wait → PostWait (код для выполнения в Go рутинах)
//   - PreStop → Stop → PostStop (остановка компонентов)
//
// Wait-функции компонентов выполняются конкурентно в отдельных горутинах.
// При ошибке в любом Wait-компоненте, возвращается первая ошибка, без graceful shutdown и выполнения последующих шагов.
func (app *App) Serve() (err error) {
	if app.IsRunning() {
		return fmt.Errorf("[compogo][%s] Server: %w", app.name, AppIsRunningErr)
	}

	app.setRunning(true)
	defer app.setRunning(false)

	app.config.StartTime = time.Now().UTC()
	defer app.closer.Close()

	app.logger.Info("Running")

	if err = app.runSteps(Init, Configuration, PreExecute, Execute, PostExecute, PreWait, Wait, PostWait); err != nil {
		return err
	}

	if err = app.wgErr.Wait(); err != nil {
		return err
	}

	app.logger.Info("Shutdown")

	return app.runSteps(PreStop, Stop, PostStop)
}

// runSteps выполняет указанные этапы для всех компонентов приложения.
func (app *App) runSteps(steps ...Step) (err error) {
	if app.parent != nil {
		if err = app.parent.runSteps(steps...); err != nil {
			return err
		}
	}

	components := app.getAllComponents()
	for _, step := range steps {
		if err = app.runStepComponents(step, components...); err != nil {
			return err
		}
	}

	return nil
}

// runStepComponents выполняет один этап для списка компонентов с учётом зависимостей.
func (app *App) runStepComponents(step Step, components ...*Component) (err error) {
	stepComponents, _ := app.steps.Get(step)

	if stopSteps.Contains(step) {
		components = slices.Clone(components)
		slices.Reverse(components)
	}

	for _, cmp := range components {
		if len(cmp.Dependencies) > 0 {
			if err = app.runStepComponents(step, cmp.Dependencies...); err != nil {
				return err
			}
		}

		if stepComponents.Contains(cmp) {
			continue
		}
		stepComponents.Add(cmp)

		if !cmp.Steps.Has(step) {
			continue
		}

		f, _ := cmp.Steps.Get(step)
		app.logger.Debugf("component '%s' executing step '%s'", cmp.Name, step)

		if step == Wait {
			app.wgErr.Go(func() error {
				return f(app.container)
			})

			continue
		}

		if err = f(app.container); err != nil {
			return err
		}
	}

	return nil
}

// getAllComponents возвращает все компоненты приложения, включая системные,
// пользовательские и компоненты родительского приложения.
func (app *App) getAllComponents() []*Component {
	var components []*Component

	if app.parent == nil {
		components = append(
			components,
			app.flagCmp,
			app.configuratorCmp,
			app.containerCmp,
			app.closerCmp,
			app.loggerCmp,
		)
	}

	if app.parent != nil {
		components = append(components, app.parent.getAllComponents()...)
	}

	components = append(components, app.components.Items()...)

	return components
}

// setRunning устанавливает флаг running для приложения и всех его родителей.
func (app *App) setRunning(isRunning bool) {
	if app.parent != nil {
		app.parent.setRunning(isRunning)
	}

	if len(app.children) > 0 {
		for _, child := range app.children {
			child.setRunning(isRunning)
		}
	}

	app.isRunning.Store(isRunning)
}

// IsRunning возвращает true, если приложение или его родитель запущены.
func (app *App) IsRunning() bool {
	if app.parent != nil {
		if isRunning := app.parent.IsRunning(); isRunning {
			return true
		}
	}

	return app.isRunning.Load()
}

// Clone создаёт дочернее приложение с указанным именем.
// Дочернее приложение наследует от родителя:
//   - DI-контейнер
//   - Конфигуратор
//   - Closer (менеджер завершения)
//   - Логгер (с вложенным именем)
//
// Используется для создания модульных приложений, где каждый модуль
// может иметь свои компоненты, но разделяет общие сервисы.
func (app *App) Clone(name string) *App {
	clone := &App{
		name:         fmt.Sprintf("%s.%s", app.name, name),
		config:       app.config,
		container:    app.container,
		configurator: app.configurator,
		closer:       app.closer,
		logger:       app.logger.GetLogger(name),
		parent:       app,
		components:   hashSlice.NewHashSlice[*Component](),
		steps:        app.steps,
	}

	app.children = append(app.children, clone)

	return clone
}
