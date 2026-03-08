package flag

import "time"

type DurationSlice interface {
	DurationSliceP(name, shorthand string, value []time.Duration, usage string) *[]time.Duration
	DurationSlice(name string, value []time.Duration, usage string) *[]time.Duration
	DurationSliceVarP(p *[]time.Duration, name, shorthand string, value []time.Duration, usage string)
	DurationSliceVar(p *[]time.Duration, name string, value []time.Duration, usage string)
	GetDurationSlice(name string) ([]time.Duration, error)
}
