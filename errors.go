package compogo

import "errors"

// Ошибки, возвращаемые при работе с приложением Compogo.
var (
	// ContainerUndefinedErr возникает, если DI-контейнер не был предоставлен через опцию WithContainer.
	// Приложение не может работать без контейнера.
	ContainerUndefinedErr = errors.New("container is undefined")

	// ConfiguratorUndefinedErr возникает, если конфигуратор не был предоставлен через опцию WithConfigurator.
	// Приложение не может загружать конфигурацию без конфигуратора.
	ConfiguratorUndefinedErr = errors.New("configurator is undefined")

	// CloserUndefinedErr возникает, если менеджер завершения (Closer) не был предоставлен через опцию WithCloser.
	// Приложение не может корректно завершиться без Closer.
	CloserUndefinedErr = errors.New("closer is undefined")

	// LoggerUndefinedErr возникает, если логгер не был предоставлен через опцию WithLogger.
	// Приложение не может логировать события без логгера.
	LoggerUndefinedErr = errors.New("logger is undefined")

	// AppIsRunningErr возникает при попытке изменить конфигурацию или компоненты
	// после запуска приложения (вызова Serve).
	AppIsRunningErr = errors.New("app is running")
)
