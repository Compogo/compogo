package closer

import (
	"context"
	"io"
)

type Closer interface {
	io.Closer
	GetContext() context.Context
	IsClosed() bool
}
