# Compogo

[https://pkg.go.dev/badge/github.com/Compogo/compogo.svg](https://pkg.go.dev/badge/github.com/Compogo/compogo.svg)[https://img.shields.io/badge/License-MIT-yellow.svg](https://img.shields.io/badge/License-MIT-yellow.svg)

Фреймворк для построения модульных Go-приложений с управлением жизненным циклом компонентов через DI-контейнер.

## Установка

```shell
go get github.com/Compogo/compogo
```

## Основные концепции

### Компонент

Компонент — единица приложения с чётко определёнными этапами жизненного цикла:

```plantuml
Init → BindFlag → Configuration
         ↓
PreExecute → Execute → PostExecute
         ↓
PreWait → Wait → PostWait
         ↓
PreStop → Stop → PostStop
```

Каждый компонент может реализовывать функции для любого из этих этапов.

#### Пример компонента:

```go
import (
	"context"
	"net/http"
    "github.com/Compogo/compogo"
)

//----- CONFIG -----

const HTTPServerAddrFieldName = "http.server.addr"

var HTTPServerAddrDefault = "localhost:8080"

type HTTPServerConfig struct {
	Addr string
}

func NewHTTPServerConfig() *HTTPServerConfig {
	return &HTTPServerConfig{}
}

func HTTPServerConfiguration(config *HTTPServerConfig, compogo.configurator Configurator) *HTTPServerConfig {
	if config.Addr == "" {
		configurator.SetDefault(HTTPServerAddrFieldName, HTTPServerAddrDefault)
		config.Addr = configurator.GetString(HTTPServerAddrFieldName)
	}
	
	return config
}

//----- CONFIG END -----

//----- SERVICE -----

type HTTPServer struct {
	config *HTTPServerConfig
	srv *http.Server
}

func NewHTTPServer(config *HTTPServerConfig) *HTTPServer {
	return &HTTPServer{
		config: config,
		srv: &http.Server{Addr: config.Addr},
	}
}

func (s *HTTPServer) Process(ctx context.Context) (err error) {
	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()
	
	go func() {
		<-ctx.Done()
		err = s.srv.Shutdown(ctx)
	}()
	
	return s.srv.ListenAndServe()
}

//----- SERVICE END -----

var HTTPServerComponent = compogo.Component{
	Name: "http_server",
	Init: compogo.StepFunc(func(container compogo.Container) error {
		return container.Provides(NewHTTPServerConfig, NewHTTPServer)
	}),
	Configuration: compogo.StepFunc(func(container compogo.Container) error {
		return container.Invoke(HTTPServerConfiguration)
	}),
	Wait: compogo.WaitFunc(func(ctx context.Context, container compogo.Container) error {
		return container.Invoke(func(httpServer *HTTPServer) error {
			return httpServer.Process(ctx)
		})
	}),
}
```

### Приложение

App управляет жизненным циклом всех компонентов:

```go
import "github.com/Compogo/compogo"

func main() {
    // 1. Создание приложения
    app := compogo.NewApp("myapp")

    // 2. Добавление компонентов
    app.AddComponents(
        dbComponent,
        httpComponent,
        metricsComponent,
    )

    // 3. Привязка флагов
    flagSet := pflag.NewFlagSet("myapp", pflag.ExitOnError)
    if err := app.BindFlags(flagSet); err != nil {
        log.Fatal(err)
    }
    flagSet.Parse(os.Args[1:])

    // 4. Запуск
    if err := app.Serve(); err != nil {
        log.Fatal(err)
    }
}
```

### Интерфейсы

#### Configurator

Загрузка конфигурации из различных источников:

```go
type Configurator interface {
    GetString(string) string
    GetInt(string) int
    GetBool(string) bool
    GetDuration(string) time.Duration
    GetStringSlice(string) []string
    // ... и другие типы
    SetDefault(string, any)
    ReadConfig() error
    With(string) Configurator
}
```

#### Container

DI-контейнер для управления зависимостями:

```go
type Container interface {
    Provide(interface{}) error
    Provides(...interface{}) error
    Invoke(interface{}) error
}
```

#### Logger

```go
type Logger interface {
    // Уровни логирования
    Panicf(string, ...interface{})
    Errorf(string, ...interface{})
    Warnf(string, ...interface{})
    Infof(string, ...interface{})
    Debugf(string, ...interface{})
    Printf(string, ...interface{})
    // Вложенные логгеры
    GetLogger(name string) Logger
}
```

#### Closer

Graceful shutdown приложения:

```go
type Closer interface {
    io.Closer
    GetContext() context.Context
    IsClosed() bool
}
```

##### Готовая реализация через сигналы ОС: 

```go
closer := compogo.NewCloserOsSignal()
go closer.Serve() // ожидание SIGINT, SIGTERM, SIGABRT
```

### Функциональные опции

```go
app := compogo.NewApp("myapp",
    // Обязательные опции
    compogo.WithContainer(container, containerCmp),
    compogo.WithConfigurator(configurator, configuratorCmp),
    compogo.WithLogger(logger, loggerCmp),
    compogo.WithCloser(closer, closerCmp),

    // Готовая опция для Closer через сигналы ОС
    compogo.WithOsSignalCloser(),

    // Добавление компонентов
    compogo.WithComponents(dbComponent, httpComponent),
)
```

### Клонирование приложений

Для модульной архитектуры:

```go
mainApp := compogo.NewApp("main")
apiModule := mainApp.Clone("api")
apiModule.AddComponents(apiComponents...)
```

### Конфигурация приложения

Автоматически добавляется компонент `compogo.Config` с полями:

```go
type Config struct {
    Cluster       string // Kubernetes cluster name
    Namespace     string // Kubernetes namespace
    ContainerId   string // Container ID
    ContainerName string // Container name
    Hostname      string // Hostname
    PID           uint64 // Process ID (автоматически)
}
```

### Метрики

Префикс для всех метрик, стандартных компонентов: `compogo_`

```go
const (
    MetricNamePrefix       = "compogo_"
    MetricAppNameFieldName = "app"
)
```

### Лицензия

```plantuml
MIT License

Copyright (c) 2026 Compogo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
