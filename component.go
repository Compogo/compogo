package compogo

import (
	"github.com/Compogo/types/linker"
)

// StepFunc — функция, выполняемая на определённом этапе жизненного цикла компонента.
// Принимает DI-контейнер для получения зависимостей.
// Возвращает ошибку, если этап не может быть завершён успешно.
type StepFunc func(container Container) error

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

	// Steps - методы обработчики жизненного цикла
	Steps *linker.Linker[Step, StepFunc]
}

type ComponentOption func(*Component)

func WithDependencies(dependencies ...*Component) ComponentOption {
	return func(c *Component) {
		c.Dependencies = append(c.Dependencies, dependencies...)
	}
}

func WithInit(f StepFunc) ComponentOption {
	return WithStep(Init, f)
}

func WithBindFlag(f StepFunc) ComponentOption {
	return WithStep(BindFlag, f)
}

func WithConfiguration(f StepFunc) ComponentOption {
	return WithStep(Configuration, f)
}

func WithPreExecute(f StepFunc) ComponentOption {
	return WithStep(PreExecute, f)
}

func WithExecute(f StepFunc) ComponentOption {
	return WithStep(Execute, f)
}

func WithPostExecute(f StepFunc) ComponentOption {
	return WithStep(PostExecute, f)
}

func WithPreWait(f StepFunc) ComponentOption {
	return WithStep(PreWait, f)
}

func WithWait(f StepFunc) ComponentOption {
	return WithStep(Wait, f)
}

func WithPostWait(f StepFunc) ComponentOption {
	return WithStep(PostWait, f)
}

func WithPreStop(f StepFunc) ComponentOption {
	return WithStep(PreStop, f)
}

func WithStop(f StepFunc) ComponentOption {
	return WithStep(Stop, f)
}

func WithPostStop(f StepFunc) ComponentOption {
	return WithStep(PostStop, f)
}

func WithStep(step Step, f StepFunc) ComponentOption {
	return func(c *Component) {
		c.Steps.Add(step, f)
	}
}

func NewComponent(name string, options ...ComponentOption) *Component {
	c := &Component{
		Name:  name,
		Steps: linker.NewLinker[Step, StepFunc](),
	}

	for _, option := range options {
		option(c)
	}

	return c
}
