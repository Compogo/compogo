package component

//go:generate stringer -type=Step

const (
	Init Step = iota

	PreRun
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
