package compogo

import "github.com/Compogo/compogo/flag"

// Option — функция настройки приложения Compogo.
// Используется для внедрения зависимостей и конфигурации при создании App через NewApp.
type Option func(app *App)

// WithLogger устанавливает логгер и компонент логгера в приложение.
// Компонент должен регистрировать логгер в DI-контейнере.
func WithLogger(logger Logger, cmp *Component) Option {
	return func(app *App) {
		app.loggerCmp = cmp
		app.logger = logger
	}
}

// WithContainer устанавливает DI-контейнер и компонент контейнера в приложение.
// Компонент должен регистрировать контейнер в DI-контейнере (самого себя).
func WithContainer(container Container, cmp *Component) Option {
	return func(app *App) {
		app.containerCmp = cmp
		app.container = container
	}
}

// WithConfigurator устанавливает конфигуратор и компонент конфигуратора в приложение.
// Компонент должен регистрировать конфигуратор в DI-контейнере.
func WithConfigurator(configurator Configurator, cmp *Component) Option {
	return func(app *App) {
		app.configuratorCmp = cmp
		app.configurator = configurator
	}
}

// WithCloser устанавливает менеджер завершения и компонент Closer в приложение.
// Компонент должен регистрировать Closer в DI-контейнере и запускать его в Wait.
func WithCloser(closer Closer, cmp *Component) Option {
	return func(app *App) {
		app.closerCmp = cmp
		app.closer = closer
	}
}

// WithComponents добавляет компоненты в приложение при создании.
// Является удобной обёрткой над AddComponents.
func WithComponents(components ...*Component) Option {
	return func(app *App) {
		_ = app.AddComponents(components...)
	}
}

// WithFlagSet создает компонент для внешнего FlagSet парсера и устанавливает его
func WithFlagSet(flagSet flag.FlagSet) Option {
	return WithFlagComponent(NewComponent(
		"flagSet.external",
		WithInit(func(container Container) error {
			return container.Provide(func() flag.FlagSet { return flagSet })
		}),
	))
}

// WithFlagComponent устанавливает компонент парсера флагов командной строки
func WithFlagComponent(cmp *Component) Option {
	return func(app *App) {
		app.flagCmp = cmp
	}
}
