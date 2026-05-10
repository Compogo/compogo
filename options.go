package compogo

import (
	"github.com/Compogo/compogo/closer"
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/configurator"
	"github.com/Compogo/compogo/container"
	"github.com/Compogo/compogo/flag"
	"github.com/Compogo/compogo/logger"
)

// Option configures the App during creation.
// Options are applied in the order they are provided to NewApp.
type Option func(app *App)

// withConfig is an internal option that automatically adds the config component.
// It ensures every app has a basic configuration structure.
func withConfig(config *Config) Option {
	return func(app *App) {
		app.config = config

		app.configCmp = &component.Component{
			Name: "compogo.Config",
			Init: component.StepFunc(func(container container.Container) error {
				return container.Provide(func() *Config { return config })
			}),
			BindFlags: component.BindFlags(func(flagSet flag.FlagSet, container container.Container) error {
				return container.Invoke(func(config *Config) {
					flagSet.StringVar(&config.Cluster, ClusterFieldName, ClusterDefault, "k8s cluster name")
					flagSet.StringVar(&config.Namespace, NamespaceFieldName, NamespaceDefault, "application namespace in k8s cluster")
					flagSet.StringVar(&config.ContainerId, ContainerIdFieldName, ContainerIdDefault, "application container id")
					flagSet.StringVar(&config.ContainerName, ContainerNameFieldName, ContainerNameDefault, "application container name")
					flagSet.StringVar(&config.Hostname, HostnameFieldName, HostnameDefault, "application hostname")
				})
			}),
			Configuration: component.StepFunc(func(container container.Container) error {
				return container.Invoke(Configuration)
			}),
		}
	}
}

// WithLogger injects a logger implementation and its component into the app.
// The logger component will participate in the application lifecycle.
func WithLogger(logger logger.Logger, cmp *component.Component) Option {
	return func(app *App) {
		app.loggerCmp = cmp
		app.logger = logger
	}
}

// WithContainer injects a DI container implementation and its component.
// The container component will participate in the application lifecycle.
func WithContainer(container container.Container, cmp *component.Component) Option {
	return func(app *App) {
		app.containerCmp = cmp
		app.container = container
	}
}

// WithConfigurator injects a configurator implementation and its component.
// The configurator component will participate in the application lifecycle.
func WithConfigurator(configurator configurator.Configurator, cmp *component.Component) Option {
	return func(app *App) {
		app.configuratorCmp = cmp
		app.configurator = configurator
	}
}

// WithCloser injects a closer implementation and its component.
// The closer component will participate in the application lifecycle.
func WithCloser(closer closer.Closer, cmp *component.Component) Option {
	return func(app *App) {
		app.closerCmp = cmp
		app.closer = closer
	}
}

// WithComponents registers additional components during app creation.
// Panics if component registration fails (should only happen on programming errors).
func WithComponents(components ...*component.Component) Option {
	return func(app *App) {
		if err := app.AddComponents(components...); err != nil {
			panic(err)
		}
	}
}
