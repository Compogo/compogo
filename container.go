package compogo

// Container определяет интерфейс DI-контейнера для управления зависимостями.
// Обеспечивает регистрацию сервисов и их получение через внедрение зависимостей.
//
// Основные операции:
//   - Provide — регистрация одного конструктора сервиса
//   - Provides — регистрация нескольких конструкторов
//   - Invoke — выполнение функции с внедрёнными зависимостями
//
// Пример использования:
//
//	container := // реализация Container
//
//	// Регистрация сервисов
//	container.Provide(func() *Config { return &Config{Port: 8080} })
//	container.Provide(func(cfg *Config) *http.Server {
//	    return &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port)}
//	})
//
//	// Получение сервиса
//	container.Invoke(func(srv *http.Server) {
//	    srv.ListenAndServe()
//	})
type Container interface {
	Provide(interface{}) error
	Provides(...interface{}) error
	Invoke(interface{}) error
}
