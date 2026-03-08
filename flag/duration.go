package flag

import "time"

type Duration interface {
	DurationP(name, shorthand string, value time.Duration, usage string) *time.Duration
	Duration(name string, value time.Duration, usage string) *time.Duration
	DurationVarP(p *time.Duration, name, shorthand string, value time.Duration, usage string)
	DurationVar(p *time.Duration, name string, value time.Duration, usage string)
	GetDuration(name string) (time.Duration, error)
}
