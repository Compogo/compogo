package compogo

// Panicer определяет интерфейс для логирования паник и критических ошибок.
// Используется для сообщений, после которых приложение должно завершиться.
type Panicer interface {
	Panicf(string, ...interface{})
	Panic(...interface{})
}

// Errorer определяет интерфейс для логирования ошибок.
// Используется для сообщений об ошибках, которые не должны приводить к панике.
type Errorer interface {
	Errorf(string, ...interface{})
	Error(...interface{})
}

// Warner определяет интерфейс для логирования предупреждений.
// Используется для сообщений о потенциальных проблемах, которые не являются ошибками.
type Warner interface {
	Warnf(string, ...interface{})
	Warn(...interface{})
}

// Informer определяет интерфейс для логирования информационных сообщений.
// Используется для сообщений о нормальном ходе выполнения приложения.
type Informer interface {
	Infof(string, ...interface{})
	Info(...interface{})
}

// Debuger определяет интерфейс для логирования отладочных сообщений.
// Используется для детальной диагностики и отладки приложения.
type Debuger interface {
	Debugf(string, ...interface{})
	Debug(...interface{})
}

// Printer определяет интерфейс для простого вывода сообщений.
// Используется для вывода без дополнительных уровней логирования.
type Printer interface {
	Printf(string, ...interface{})
	Print(...interface{})
}

// Logger объединяет все интерфейсы логирования и предоставляет
// возможность создания вложенных логгеров с префиксами.
//
// Logger реализует уровни логирования:
//   - Panic: критическая ошибка, вызывающая панику
//   - Error: ошибка, требующая внимания
//   - Warn: предупреждение о потенциальной проблеме
//   - Info: информационное сообщение
//   - Debug: отладочная информация
//   - Print: простой вывод
//
// Пример использования:
//
//	logger := // реализация Logger
//	logger.Info("Application started")
//	logger.WithField("port", 8080).Info("Server listening")
//
//	// Вложенные логгеры для модулей
//	httpLogger := logger.GetLogger("http")
//	httpLogger.Info("HTTP server initialized")
type Logger interface {
	Panicer
	Errorer
	Warner
	Informer
	Debuger
	Printer

	// GetLogger возвращает вложенный логгер с указанным именем.
	// Имя используется для добавления префикса или контекста к сообщениям.
	// Позволяет создавать иерархию логгеров для разных компонентов.
	//
	// Пример:
	//
	//	appLogger := logger.GetLogger("app")
	//	dbLogger := appLogger.GetLogger("database")
	//	dbLogger.Info("Connected to database")
	GetLogger(name string) Logger
}
