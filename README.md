# Compogo ⚛️

**Compogo** — это легковесный и прагматичный способ сборки Go-сервисов из декларативных компонентов. Забудьте про многостраничный `main.go` — просто опишите, от чего зависит ваш компонент, что он делает на каждом этапе жизни, и доверьте ядру сборку всего приложения в цельную, работающую систему.

## 🚀 Быстрый старт

```go
package main

import (
    "github.com/Compogo/compogo"
    "github.com/Compogo/compogo/logger"
    "github.com/Compogo/compogo/container"
    "github.com/Compogo/compogo/configurator"
    "github.com/Compogo/myapp/http"
    "github.com/Compogo/myapp/service"
)

func main() {
    app := compogo.NewApp("myapp",
        compogo.WithOsSignalCloser(),           // graceful shutdown по Ctrl+C
        compogo.WithLogger(logger.NewSlog(), logger.Component),
        compogo.WithContainer(container.NewDig(), container.Component),
        compogo.WithConfigurator(configurator.NewViper(), configurator.Component),
        compogo.WithComponents(
            http.Component,
            service.Component,
        ),
    )

    if err := app.Serve(); err != nil {
        panic(err)
    }
}
```

## ✨ Ключевые возможности

### 🧩 Декларативные компоненты

Каждый компонент сам описывает свои зависимости и поведение на каждом этапе жизни:

```go
var Component = &component.Component{
    Dependencies: app.Components{logger.Component, config.Component},
    
    Init: func(c container.Container) error {
        return c.Provide(NewMyService)
    },
    
    PostRun: func(c container.Container) error {
        return c.Invoke(func(s *MyService, r runner.Runner) {
            r.RunProcess(s)
        })
    },
}
```

### 📋 9-шаговый жизненный цикл

`PreRun → Run → PostRun → PreWait → Wait → PostWait → PreStop → Stop → PostStop`

### ⏱️ Таймауты на каждый шаг

Каждый этап имеет свой таймаут — никакой шаг не зависнет навсегда.

### 🔌 Graceful shutdown из коробки

Одна строка — и ваше приложение корректно завершается по SIGINT/SIGTERM.

### 📦 Установка

```bash
go get github.com/Compogo/compogo
```

### 📄 Лицензия

Apache 2.0 — смотрите файл [LICENSE](LICENSE).

Compogo — создан с ❤️ для Go-разработчиков, которые ценят порядок и предсказуемость.
