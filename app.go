package compogo

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Compogo/compogo/flag"
	hashSlice "github.com/Compogo/types/hash_slice"
	"github.com/Compogo/types/linker"
	"github.com/Compogo/types/set"
)

// App — главный объект приложения Compogo.
// Управляет жизненным циклом всех компонентов, обеспечивает:
//   - Инициализацию и конфигурацию компонентов в правильном порядке (с учётом зависимостей)
//   - Привязку флагов командной строки
//   - Graceful shutdown через контекст и WaitGroup
//   - Иерархию приложений (под-приложения через Clone)
//
// App потокобезопасен и может использоваться из нескольких горутин.
type App struct {
	name      string
	startTime time.Time

	appComponent *Component

	configCmp *Component
	config    *Config

	containerCmp *Component
	container    Container

	configuratorCmp *Component
	configurator    Configurator

	closerCmp *Component
	closer    Closer

	loggerCmp *Component
	logger    Logger

	components *hashSlice.HashSlice[*Component]
	wg         sync.WaitGroup
	waitMutex  sync.Mutex

	steps *linker.Linker[Step, set.Set[*Component]]

	isRunning atomic.Bool

	parent *App
}

// NewApp создаёт новое приложение Compogo с указанным именем.
// Имя используется в логах и сообщениях об ошибках.
//
// Опции могут быть переданы для конфигурации приложения:
//   - WithConfigurator — установка загрузчика конфигурации
//   - WithLogger — установка логгера
//   - WithContainer — установка DI-контейнера
//   - WithCloser — установка менеджера завершения
//
// Пример:
//
//	app := NewApp("myapp",
//	    WithLogger(logrusLogger),
//	    WithConfigurator(viperConfigurator),
//	)
func NewApp(name string, options ...Option) *App {
	app := &App{
		name: name,
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
	}

	options = append(options, withConfig(NewConfig()))

	for _, option := range options {
		option(app)
	}

	app.appComponent = &Component{
		Init: func(container Container) error {
			return container.Provide(func() *App {
				return app
			})
		},
	}

	return app
}

// AddComponents добавляет компоненты в приложение.
// Компоненты добавляются рекурсивно вместе со всеми зависимостями.
// Дубликаты автоматически исключаются.
//
// Важно: компоненты нельзя добавлять после запуска приложения (вызова Serve).
// В этом случае возвращается ошибка AppIsRunningError.
//
// Пример:
//
//	app := NewApp("myapp")
//	err := app.AddComponents(
//	    &component.Component{Name: "db", Init: initDB},
//	    &component.Component{Name: "api", Dependencies: component.Components{dbComponent}},
//	)
func (app *App) AddComponents(components ...*Component) (err error) {
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

// existComponent проверяет существование компонента в приложении или его родителе.
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
func (app *App) BindFlags(flagSet flag.FlagSet) (err error) {
	if err = app.validate(); err != nil {
		return err
	}

	if app.IsRunning() {
		return fmt.Errorf("[compogo][%s] BindFlags: %w", app.name, AppIsRunningError)
	}

	if err = app.runComponents(Init); err != nil {
		return fmt.Errorf("[compogo][%s] BindFlags: %w", app.name, err)
	}

	if app.parent != nil {
		if err = app.parent.BindFlags(flagSet); err != nil {
			return err
		}
	}

	components := app.getCoreComponents()
	components = append(components, app.components.Items()...)

	bindFlags, _ := app.steps.Get(BindFlag)

	return app.bindFlags(flagSet, bindFlags, components...)
}

// bindFlags — внутренняя рекурсивная реализация привязки флагов.
func (app *App) bindFlags(flagSet flag.FlagSet, bindFlags set.Set[*Component], components ...*Component) (err error) {
	for _, cmp := range components {
		if len(cmp.Dependencies) > 0 {
			if err = app.bindFlags(flagSet, bindFlags, cmp.Dependencies...); err != nil {
				return err
			}
		}

		if !bindFlags.Contains(cmp) && cmp.BindFlags != nil {
			if err = cmp.BindFlags(flagSet, app.container); err != nil {
				return err
			}

			bindFlags.Add(cmp)
		}
	}

	return nil
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
// При получении сигнала завершения (через Closer) или ошибке в любом Wait-компоненте,
// инициируется graceful shutdown.
func (app *App) Serve() (err error) {
	if err = app.validate(); err != nil {
		return err
	}

	if app.IsRunning() {
		return fmt.Errorf("[compogo][%s] %w", app.name, AppIsRunningError)
	}

	app.setRunning(true)
	defer app.setRunning(false)

	app.startTime = time.Now().UTC()
	l := app.logger.GetLogger("compogo").GetLogger(app.name)

	l.Info("Running")

	if err = app.runComponents(Init, Configuration, PreExecute, Execute, PostExecute, PreWait); err != nil {
		return err
	}

	components := app.getAllComponents()

	chainErr := make(chan error, 1)
	defer close(chainErr)

	ctx, cancelFunc := context.WithCancel(app.closer.GetContext())
	defer cancelFunc()

	waitComponents, _ := app.steps.Get(Wait)

	for _, cmp := range components {
		if cmp.Wait != nil && !waitComponents.Contains(cmp) {
			app.serveComponent(ctx, cmp, waitComponents, chainErr)
			waitComponents.Add(cmp)
		}
	}

	if err = app.runComponents(PostWait); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				l.Info("Shutdown")
				return
			case err := <-chainErr:
				l.Error(err.Error())
				if err := app.closer.Close(); err != nil {
					l.Error(err.Error())
				}
			}
		}
	}()

	app.wg.Wait()

	if err = app.runComponents(PreStop, Stop, PostStop); err != nil {
		return err
	}

	return err
}

// runComponents выполняет указанные этапы для всех компонентов приложения.
func (app *App) runComponents(steps ...Step) (err error) {
	if app.parent != nil {
		if err = app.parent.runComponents(steps...); err != nil {
			return err
		}
	}

	components := app.getCoreComponents()
	components = append(components, app.components.Items()...)

	for _, step := range steps {
		if err = app.runStepComponents(step, components...); err != nil {
			return err
		}
	}

	return nil
}

// runStepComponents выполняет один этап для списка компонентов с учётом зависимостей.
func (app *App) runStepComponents(step Step, components ...*Component) (err error) {
	var fnc StepFunc
	stepComponents, _ := app.steps.Get(step)

	if stopSteps.Contains(step) {
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

		fnc, err = cmp.stepFunc(step)
		if err != nil {
			return fmt.Errorf("[compogo][%s][%s]: %w", app.name, step, err)
		}

		if fnc == nil {
			continue
		}

		if err := fnc(app.container); err != nil {
			return err
		}

		stepComponents.Add(cmp)
	}

	return nil
}

// serveComponent запускает Wait-функцию компонента в отдельной горутине.
// Использует sync.WaitGroup для отслеживания завершения всех Wait-функций.
// При панике или ошибке отправляет сообщение в канал chainErr.
func (app *App) serveComponent(ctx context.Context, cmp *Component, waitComponents set.Set[*Component], chainErr chan error) {
	app.waitMutex.Lock()
	defer app.waitMutex.Unlock()

	app.wg.Go(func(ctx context.Context, cmp *Component, waitComponents set.Set[*Component]) func() {
		return func() {
			defer func() {
				app.waitMutex.Lock()
				defer app.waitMutex.Unlock()
				waitComponents.Remove(cmp)
			}()

			defer func() {
				if r := recover(); r != nil {
					chainErr <- fmt.Errorf("component '%s' wait failed: %s", cmp.Name, r)
				}
			}()

			ctx, cancelFunc := context.WithCancel(ctx)
			defer cancelFunc()

			if err := cmp.Wait(ctx, app.container); err != nil {
				chainErr <- fmt.Errorf("component '%s' wait failed: %w", cmp.Name, err)
			}
		}
	}(ctx, cmp, waitComponents))
}

// getAllComponents возвращает все компоненты приложения, включая системные,
// пользовательские и компоненты родительского приложения.
func (app *App) getAllComponents() []*Component {
	components := app.getCoreComponents()

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

// validate проверяет, что все необходимые зависимости приложения инициализированы.
// Возвращает агрегированную ошибку, если чего-то не хватает.
func (app *App) validate() error {
	var err error

	if app.parent != nil {
		return app.parent.validate()
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
	return &App{
		name:         fmt.Sprintf("%s.%s", app.name, name),
		appComponent: app.appComponent,
		config:       app.config,
		container:    app.container,
		configurator: app.configurator,
		closer:       app.closer,
		logger:       app.logger.GetLogger("compogo").GetLogger(app.name).GetLogger(name),
		parent:       app,
		components:   hashSlice.NewHashSlice[*Component](),
		steps:        app.steps,
	}
}

// getCoreComponents возвращает системные компоненты приложения.
// Эти компоненты всегда выполняются первыми.
func (app *App) getCoreComponents() []*Component {
	if app.configCmp == nil {
		return nil
	}

	return []*Component{
		app.configuratorCmp,
		app.containerCmp,
		app.configCmp,
		app.closerCmp,
		app.loggerCmp,
		app.appComponent,
	}
}

// StartTime возвращает время запуска приложения.
// Полезно для метрик и мониторинга uptime.
func (app *App) StartTime() time.Time {
	return app.startTime
}
