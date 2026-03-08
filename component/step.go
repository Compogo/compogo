package component

//go:generate stringer -type=Step

const (
	PreRun Step = iota
	Run
	PostRun

	PreWait
	Wait
	PostWait

	PreStop
	Stop
	PostStop
)

type Step uint8
