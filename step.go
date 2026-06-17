package compogo

import "github.com/Compogo/types/set"

//go:generate stringer -type=Step

// Step определяет этапы жизненного цикла компонента.
// Каждый компонент может реализовать функции для любого из этих этапов.
//
// Жизненный цикл компонента:
//
//	Init → BindFlag → Configuration
//	         ↓
//	PreExecute → Execute → PostExecute
//	         ↓
//	PreWait → Wait → PostWait
//	         ↓
//	PreStop → Stop → PostStop
//
// Где:
//   - Init: инициализация компонента, регистрация зависимостей в DI-контейнере
//   - BindFlag: привязка флагов командной строки
//   - Configuration: загрузка конфигурации
//   - PreExecute/Execute/PostExecute: основной рабочий цикл
//   - PreWait/Wait/PostWait: ожидание сигналов завершения
//   - PreStop/Stop/PostStop: graceful shutdown
const (
	// Init — начальный этап инициализации компонента.
	// На этом этапе компонент регистрирует свои зависимости и сервисы в DI-контейнере.
	Init Step = iota

	// BindFlag — этап привязки флагов командной строки.
	// Компонент определяет свои флаги и связывает их с полями конфигурации.
	BindFlag

	// Configuration — этап загрузки конфигурации.
	// Компонент читает и валидирует свою конфигурацию.
	Configuration

	// PreExecute — подготовка к выполнению.
	// Выполняется перед основным рабочим циклом компонента.
	PreExecute

	// Execute — основной рабочий цикл компонента.
	// Обычно содержит блокирующую логику (например, запуск HTTP-сервера).
	Execute

	// PostExecute — завершение рабочего цикла.
	// Выполняется после завершения Execute.
	PostExecute

	// PreWait — подготовка к ожиданию сигналов.
	PreWait

	// Wait — ожидание сигналов завершения (например, SIGINT, SIGTERM).
	// Компонент может ожидать свой сигнал или использовать общий контекст приложения.
	Wait

	// PostWait — завершение ожидания.
	PostWait

	// PreStop — подготовка к остановке.
	PreStop

	// Stop — остановка компонента.
	// Здесь компонент должен корректно завершить все рабочие горутины.
	Stop

	// PostStop — завершение остановки, финальная очистка.
	PostStop
)

// Step представляет этап жизненного цикла компонента.
type Step uint8

var stopSteps = set.NewSet[Step](PreStop, Stop, PostStop)
