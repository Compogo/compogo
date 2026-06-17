package compogo

// Константы для метрик приложения Compogo.
const (
	// MetricNamePrefix — префикс для всех метрик, собираемых Compogo.
	// Используется для идентификации метрик фреймворка в системах мониторинга.
	//
	// Пример:
	//   compogo_app_start_time
	//   compogo_component_execute_duration
	MetricNamePrefix = "compogo_"

	// MetricAppNameFieldName — имя поля в метриках, содержащее имя приложения.
	// Используется как label для фильтрации метрик по приложению.
	//
	// Пример использования в метрике:
	//   compogo_app_info{app="myapp", version="1.0.0"}
	MetricAppNameFieldName = "app"
)
