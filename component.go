package compogo

import (
	"context"
	"errors"
	"fmt"

	"github.com/Compogo/compogo/flag"
)

// Ошибки, возвращаемые при работе с компонентами.
var (
	// StepUndefinedError возникает при попытке получить функцию для несуществующего этапа.
	StepUndefinedError = errors.New("step is undefined")
)

// StepFunc — функция, выполняемая на определённом этапе жизненного цикла компонента.
// Принимает DI-контейнер для получения зависимостей.
// Возвращает ошибку, если этап не может быть завершён успешно.
type StepFunc func(container Container) error

// WaitFunc — специальная функция для этапа Wait.
// Отличается от StepFunc наличием контекста, который может быть отменён при завершении приложения.
// Позволяет компоненту корректно завершить работу по сигналу.
//
// Пример:
//
//	func wait(ctx context.Context, c container.Container) error {
//	    <-ctx.Done() // ожидание сигнала завершения
//	    return nil
//	}
type WaitFunc func(ctx context.Context, container Container) error

// BindFlags — функция для привязки флагов командной строки.
// Получает FlagSet для регистрации флагов и DI-контейнер для доступа к конфигурации.
type BindFlags func(flagSet flag.FlagSet, container Container) error

// Components — слайс компонентов для удобного группового добавления.
type Components []*Component

// Component представляет собой единицу приложения Compogo.
// Каждый компонент имеет имя, может зависеть от других компонентов
// и реализовывать функции для различных этапов жизненного цикла.
//
// Компоненты должны быть декларативными — описывать ЧТО делать, а не КАК.
// DI-контейнер предоставляет все необходимые зависимости.
type Component struct {
	// Name — уникальное имя компонента. Используется в логах и сообщениях об ошибках.
	Name string

	// Dependencies — компоненты, от которых зависит текущий компонент.
	// Compogo автоматически обеспечивает правильный порядок запуска и остановки.
	// Если компонент A зависит от B, то B будет запущен раньше A и остановлен позже.
	Dependencies Components

	// Init — инициализация компонента.
	// Здесь следует регистрировать сервисы в DI-контейнере через container.Provide().
	Init StepFunc

	// BindFlags — привязка флагов командной строки.
	// Используйте flagSet для регистрации флагов компонента.
	BindFlags BindFlags

	// Configuration — загрузка конфигурации.
	Configuration StepFunc

	// PreExecute — подготовка к выполнению.
	PreExecute StepFunc

	// Execute — выполнение рабочего кода в однопоточном режиме.
	// Здесь не должно быть долгих блокирующих вызовов по типу ListenAndServe().
	Execute StepFunc

	// PostExecute — завершение после однопоточного рабочего кода.
	PostExecute StepFunc

	// PreWait — подготовка к ожиданию.
	PreWait StepFunc

	// Wait — код который должен быть выполнен в отдельной GO рутине.
	// Получает контекст, который отменяется при завершении приложения.
	// Здесь уже можно вызывать блокирующие вызовы по типу ListenAndServe().
	Wait WaitFunc

	// PostWait — код который должен выполнится после выполнения блокирующего вызова.
	PostWait StepFunc

	// PreStop — подготовка к остановке.
	PreStop StepFunc

	// Stop — остановка компонента.
	// Здесь следует завершить все рабочие горутины и освободить ресурсы.
	Stop StepFunc

	// PostStop — завершение остановки, финальная очистка.
	PostStop StepFunc
}

// stepFunc возвращает функцию, соответствующую указанному этапу жизненного цикла.
// Если этап не поддерживается компонентом, возвращает StepUndefinedError.
func (c *Component) stepFunc(step Step) (StepFunc, error) {
	switch step {
	case Init:
		return c.Init, nil
	case Configuration:
		return c.Configuration, nil
	case PreExecute:
		return c.PreExecute, nil
	case Execute:
		return c.Execute, nil
	case PostExecute:
		return c.PostExecute, nil
	case PreWait:
		return c.PreWait, nil
	case PostWait:
		return c.PostWait, nil
	case PreStop:
		return c.PreStop, nil
	case Stop:
		return c.Stop, nil
	case PostStop:
		return c.PostStop, nil
	default:
		return nil, fmt.Errorf("[component.%s] step '%s' %w", c.Name, step, StepUndefinedError)
	}
}
