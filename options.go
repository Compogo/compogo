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

					flagSet.DurationVar(&config.InitDuration, InitDurationFieldName, InitDurationDefault, "maximum time to wait for a component's response at the init step")

					flagSet.DurationVar(&config.PreRunDuration, PreRunDurationFieldName, PreRunDurationDefault, "maximum time to wait for a component's response at the pre-run step")
					flagSet.DurationVar(&config.RunDuration, RunDurationFieldName, RunDurationDefault, "maximum time to wait for a component's response at the run step")
					flagSet.DurationVar(&config.PostRunDuration, PostRunDurationFieldName, PostRunDurationDefault, "maximum time to wait for a component's response at the post-run step")

					flagSet.DurationVar(&config.PreWaitDuration, PreWaitDurationFieldName, PreWaitDurationDefault, "maximum time to wait for a component's response at the pre-wait step")
					flagSet.DurationVar(&config.PostWaitDuration, PostWaitDurationFieldName, PostWaitDurationDefault, "maximum time to wait for a component's response at the post-wait step")

					flagSet.DurationVar(&config.PreStopDuration, PreStopDurationFieldName, PreStopDurationDefault, "maximum time to wait for a component's response at the pre-stop step")
					flagSet.DurationVar(&config.StopDuration, StopDurationFieldName, StopDurationDefault, "maximum time to wait for a component's response at the stop step")
					flagSet.DurationVar(&config.PostStopDuration, PostStopDurationFieldName, PostStopDurationDefault, "maximum time to wait for a component's response at the post-stop step")
				})
			}),
			PreRun: component.StepFunc(func(container container.Container) error {
				if err := container.Invoke(Configuration); err != nil {
					return err
				}

				return container.Invoke(func(config *Config) {
					config.Name = app.name
				})
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
