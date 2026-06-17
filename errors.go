package compogo

import "errors"

// Ошибки, возвращаемые при работе с приложением Compogo.
var (
	// ContainerUndefinedError возникает, если DI-контейнер не был предоставлен через опцию WithContainer.
	// Приложение не может работать без контейнера.
	ContainerUndefinedError = errors.New("container is undefined")

	// ConfiguratorUndefinedError возникает, если конфигуратор не был предоставлен через опцию WithConfigurator.
	// Приложение не может загружать конфигурацию без конфигуратора.
	ConfiguratorUndefinedError = errors.New("configurator is undefined")

	// CloserUndefinedError возникает, если менеджер завершения (Closer) не был предоставлен через опцию WithCloser.
	// Приложение не может корректно завершиться без Closer.
	CloserUndefinedError = errors.New("closer is undefined")

	// LoggerUndefinedError возникает, если логгер не был предоставлен через опцию WithLogger.
	// Приложение не может логировать события без логгера.
	LoggerUndefinedError = errors.New("logger is undefined")

	// AppIsRunningError возникает при попытке изменить конфигурацию или компоненты
	// после запуска приложения (вызова Serve).
	AppIsRunningError = errors.New("app is running")
)
