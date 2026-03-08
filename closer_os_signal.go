package compogo

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"

	closer2 "github.com/Compogo/compogo/closer"
	"github.com/Compogo/compogo/component"
	"github.com/Compogo/compogo/container"
)

type CloserOsSignal struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	exitSignal chan os.Signal
	isClosed   atomic.Bool
}

func NewCloserOsSignal() *CloserOsSignal {
	closer := &CloserOsSignal{}

	closer.ctx, closer.cancelFunc = context.WithCancel(context.Background())

	return closer
}

func (closer *CloserOsSignal) Close() error {
	closer.cancelFunc()
	closer.isClosed.Store(true)
	return nil
}

func (closer *CloserOsSignal) GetContext() context.Context {
	return closer.ctx
}

func (closer *CloserOsSignal) IsClosed() bool {
	return closer.isClosed.Load()
}

func (closer *CloserOsSignal) Serve() {
	signal.Notify(closer.exitSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)

	<-closer.exitSignal

	_ = closer.Close()
}

// WithOsSignalCloser returns an Option that adds a default closer listening for OS signals.
// The closer handles SIGINT, SIGTERM, and SIGABRT, providing graceful shutdown out of the box.
// Perfect for standard applications - just add this option and your app stops gracefully on Ctrl+C.
func WithOsSignalCloser() Option {
	closer := NewCloserOsSignal()

	return WithCloser(closer, &component.Component{
		Init: component.InitFunc(func(container container.Container) error {
			return container.Provides(
				func() *CloserOsSignal { return closer },
				func(cl *CloserOsSignal) closer2.Closer { return cl },
			)
		}),
		Wait: component.StepFunc(func(container container.Container) error {
			return container.Invoke(func(cl *CloserOsSignal) {
				cl.Serve()
			})
		}),
	})
}
