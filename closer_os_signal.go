package compogo

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

// CloserOsSignal реализует интерфейс Closer для graceful shutdown
// по сигналам операционной системы (SIGINT, SIGTERM, SIGABRT).
//
// При получении одного из этих сигналов вызывается Close(), который
// отменяет контекст и переводит Closer в закрытое состояние.
type CloserOsSignal struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	exitSignal chan os.Signal
	isClosed   atomic.Bool
}

// NewCloserOsSignal создаёт новый CloserOsSignal.
// Контекст создаётся через context.WithCancel, канал для сигналов — с буфером 1.
//
// Пример:
//
//	closer := NewCloserOsSignal()
//	go closer.Serve() // запускаем ожидание сигналов
//	<-closer.GetContext().Done() // блокируемся до получения сигнала
func NewCloserOsSignal() *CloserOsSignal {
	closer := &CloserOsSignal{
		exitSignal: make(chan os.Signal, 1),
	}

	closer.ctx, closer.cancelFunc = context.WithCancel(context.Background())

	return closer
}

// Close отменяет контекст и устанавливает флаг isClosed в true.
// Реализует интерфейс io.Closer.
func (closer *CloserOsSignal) Close() error {
	closer.cancelFunc()
	closer.isClosed.Store(true)
	return nil
}

// GetContext возвращает контекст, который отменяется при вызове Close().
func (closer *CloserOsSignal) GetContext() context.Context {
	return closer.ctx
}

// IsClosed возвращает true, если Close() уже был вызван.
func (closer *CloserOsSignal) IsClosed() bool {
	return closer.isClosed.Load()
}

// Serve запускает ожидание сигналов ОС.
// При получении SIGINT, SIGTERM или SIGABRT вызывает Close().
// Обычно запускается в отдельной горутине.
//
// Пример:
//
//	closer := NewCloserOsSignal()
//	go closer.Serve()
//	<-closer.GetContext().Done()
//	// приложение завершается
func (closer *CloserOsSignal) Serve() {
	signal.Notify(closer.exitSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	<-closer.exitSignal

	_ = closer.Close()
}

// WithOsSignalCloser возвращает Option для добавления CloserOsSignal в приложение.
// Опция:
//   - Создаёт новый CloserOsSignal
//   - Устанавливает его как Closer через WithCloser
//   - Добавляет компонент, который регистрирует CloserOsSignal в DI-контейнере
//   - В Wait-функции компонента запускает Serve() для ожидания сигналов ОС
//
// Пример использования в приложении:
//
//	app := NewApp("myapp",
//	    WithOsSignalCloser(),
//	    // другие опции...
//	)
func WithOsSignalCloser() Option {
	closer := NewCloserOsSignal()

	return WithCloser(closer, &Component{
		Name: "closer.OsSignal",
		Init: StepFunc(func(container Container) error {
			return container.Provides(
				func() *CloserOsSignal { return closer },
				func(cl *CloserOsSignal) Closer { return cl },
			)
		}),
		Wait: WaitFunc(func(_ context.Context, container Container) error {
			return container.Invoke(func(cl *CloserOsSignal) {
				cl.Serve()
			})
		}),
	})
}
