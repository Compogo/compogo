package component

//go:generate stringer -type=Step

const (
	Init Step = iota
	Configuration

	PreExecute
	Execute
	PostExecute

	PreWait
	Wait
	PostWait

	PreStop
	Stop
	PostStop
)

type Step uint8
