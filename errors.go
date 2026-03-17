package compogo

import "errors"

var (
	ContainerUndefinedError    = errors.New("container is undefined")
	ConfiguratorUndefinedError = errors.New("configurator is undefined")
	CloserUndefinedError       = errors.New("closer is undefined")
	LoggerUndefinedError       = errors.New("logger is undefined")
	AppIsRunningError          = errors.New("app is running")
	ComponentStepTimeoutError  = errors.New("timeout execution of the component at the step")
)
