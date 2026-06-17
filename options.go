package compogo

import (
	"os"

	"github.com/Compogo/compogo/flag"
)

// Option — функция настройки приложения Compogo.
// Используется для внедрения зависимостей и конфигурации при создании App через NewApp.
//
// Пример:
//
//	app := NewApp("myapp",
//	    WithLogger(myLogger, loggerComponent),
//	    WithContainer(myContainer, containerComponent),
//	    WithConfigurator(myConfigurator, configuratorComponent),
//	)
type Option func(app *App)

// withConfig — внутренняя опция для установки конфигурации приложения.
// Создаёт системный компонент compogo.Config, который:
//   - Регистрирует Config в DI-контейнере
//   - Привязывает флаги командной строки для кластерных настроек
//   - Загружает конфигурацию из Configurator
//   - Устанавливает значения по умолчанию для полей конфигурации
//
// Поля конфигурации:
//   - Cluster: имя Kubernetes кластера
//   - Namespace: пространство имён в Kubernetes
//   - ContainerId: ID контейнера
//   - ContainerName: имя контейнера
//   - Hostname: имя хоста
//   - PID: идентификатор процесса (устанавливается автоматически)
//
// Эта опция применяется автоматически в NewApp и не должна использоваться пользователем.
func withConfig(config *Config) Option {
	return func(app *App) {
		app.config = config

		app.configCmp = &Component{
			Name: "compogo.Config",
			Init: StepFunc(func(container Container) error {
				return container.Provide(func() *Config { return config })
			}),
			BindFlags: BindFlags(func(flagSet flag.FlagSet, container Container) error {
				return container.Invoke(func(config *Config) {
					flagSet.StringVar(&config.Cluster, ClusterFieldName, ClusterDefault, "k8s cluster name")
					flagSet.StringVar(&config.Namespace, NamespaceFieldName, NamespaceDefault, "application namespace in k8s cluster")
					flagSet.StringVar(&config.ContainerId, ContainerIdFieldName, ContainerIdDefault, "application container id")
					flagSet.StringVar(&config.ContainerName, ContainerNameFieldName, ContainerNameDefault, "application container name")
					flagSet.StringVar(&config.Hostname, HostnameFieldName, HostnameDefault, "application hostname")
				})
			}),
			Configuration: StepFunc(func(container Container) error {
				return container.Invoke(func(config *Config, configurator Configurator) *Config {
					config.PID = uint64(os.Getpid())

					if config.Cluster == "" || config.Cluster == ClusterDefault {
						configurator.SetDefault(ClusterFieldName, ClusterDefault)
						config.Cluster = configurator.GetString(ClusterFieldName)
					}

					if config.Namespace == "" || config.Namespace == NamespaceDefault {
						configurator.SetDefault(NamespaceFieldName, NamespaceDefault)
						config.Namespace = configurator.GetString(NamespaceFieldName)
					}

					if config.ContainerId == "" || config.ContainerId == ContainerIdDefault {
						configurator.SetDefault(ContainerIdFieldName, ContainerIdDefault)
						config.ContainerId = configurator.GetString(ContainerIdFieldName)
					}

					if config.ContainerName == "" || config.ContainerName == ContainerNameDefault {
						configurator.SetDefault(ContainerNameFieldName, ContainerNameDefault)
						config.ContainerName = configurator.GetString(ContainerNameFieldName)
					}

					if config.Hostname == "" || config.Hostname == HostnameDefault {
						configurator.SetDefault(HostnameFieldName, HostnameDefault)
						config.Hostname = configurator.GetString(HostnameFieldName)
					}

					return config
				})
			}),
		}
	}
}

// WithLogger устанавливает логгер и компонент логгера в приложение.
// Компонент должен регистрировать логгер в DI-контейнере.
//
// Пример:
//
//	logger := logrus.New()
//	loggerCmp := &Component{
//	    Name: "logger",
//	    Init: func(c Container) error {
//	        return c.Provide(func() Logger { return logger })
//	    },
//	}
//	app := NewApp("myapp", WithLogger(logger, loggerCmp))
func WithLogger(logger Logger, cmp *Component) Option {
	return func(app *App) {
		app.loggerCmp = cmp
		app.logger = logger
	}
}

// WithContainer устанавливает DI-контейнер и компонент контейнера в приложение.
// Компонент должен регистрировать контейнер в DI-контейнере (самого себя).
//
// Пример:
//
//	container := fx.New()
//	containerCmp := &Component{
//	    Name: "container",
//	    Init: func(c Container) error {
//	        return c.Provide(func() Container { return container })
//	    },
//	}
//	app := NewApp("myapp", WithContainer(container, containerCmp))
func WithContainer(container Container, cmp *Component) Option {
	return func(app *App) {
		app.containerCmp = cmp
		app.container = container
	}
}

// WithConfigurator устанавливает конфигуратор и компонент конфигуратора в приложение.
// Компонент должен регистрировать конфигуратор в DI-контейнере.
//
// Пример:
//
//	configurator := viper.New()
//	configuratorCmp := &Component{
//	    Name: "configurator",
//	    Init: func(c Container) error {
//	        return c.Provide(func() Configurator { return configurator })
//	    },
//	}
//	app := NewApp("myapp", WithConfigurator(configurator, configuratorCmp))
func WithConfigurator(configurator Configurator, cmp *Component) Option {
	return func(app *App) {
		app.configuratorCmp = cmp
		app.configurator = configurator
	}
}

// WithCloser устанавливает менеджер завершения и компонент Closer в приложение.
// Компонент должен регистрировать Closer в DI-контейнере и запускать его в Wait.
//
// Пример:
//
//	closer := NewCloserOsSignal()
//	closerCmp := &Component{
//	    Name: "closer",
//	    Init: func(c Container) error {
//	        return c.Provide(func() Closer { return closer })
//	    },
//	    Wait: func(ctx context.Context, c Container) error {
//	        return c.Invoke(func(cl *CloserOsSignal) {
//	            cl.Serve()
//	        })
//	    },
//	}
//	app := NewApp("myapp", WithCloser(closer, closerCmp))
func WithCloser(closer Closer, cmp *Component) Option {
	return func(app *App) {
		app.closerCmp = cmp
		app.closer = closer
	}
}

// WithComponents добавляет компоненты в приложение при создании.
// Является удобной обёрткой над AddComponents.
// Если добавление компонентов вызывает ошибку, паникует.
//
// Пример:
//
//	app := NewApp("myapp",
//	    WithComponents(
//	        dbComponent,
//	        httpComponent,
//	        metricsComponent,
//	    ),
//	)
func WithComponents(components ...*Component) Option {
	return func(app *App) {
		if err := app.AddComponents(components...); err != nil {
			panic(err)
		}
	}
}
