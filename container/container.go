package container

type Container interface {
	Provide(interface{}) error
	Provides(...interface{}) error
	Invoke(interface{}) error
}
