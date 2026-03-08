package container

type Container interface {
	Provide(any) error
	Provides(...any) error
	Invoke(any) error
}
