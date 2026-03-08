package logger

type Panicer interface {
	Panicf(string, ...interface{})
	Panic(...interface{})
}

type Errorer interface {
	Errorf(string, ...interface{})
	Error(...interface{})
}

type Warner interface {
	Warnf(string, ...interface{})
	Warn(...interface{})
}

type Informer interface {
	Infof(string, ...interface{})
	Info(...interface{})
}

type Debuger interface {
	Debugf(string, ...interface{})
	Debug(...interface{})
}

type Printer interface {
	Printf(string, ...interface{})
	Print(...interface{})
}

// Logger : General slog interface for any implementation.
type Logger interface {
	Panicer
	Errorer
	Warner
	Informer
	Debuger
	Printer

	GetLogger(name string) Logger
}
