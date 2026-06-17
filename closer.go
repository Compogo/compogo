package compogo

import (
	"context"
	"io"
)

// Closer объединяет интерфейс io.Closer с возможностью получения контекста
// и проверки состояния закрытия.
//
// Используется для graceful shutdown приложения:
//   - Close() инициирует завершение работы
//   - GetContext() возвращает контекст, который отменяется при вызове Close()
//   - IsClosed() позволяет проверить, был ли уже вызван Close()
type Closer interface {
	io.Closer
	GetContext() context.Context
	IsClosed() bool
}
