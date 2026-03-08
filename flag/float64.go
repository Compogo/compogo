package flag

type Float64 interface {
	Float64P(name, shorthand string, value float64, usage string) *float64
	Float64(name string, value float64, usage string) *float64
	Float64VarP(p *float64, name, shorthand string, value float64, usage string)
	Float64Var(p *float64, name string, value float64, usage string)
	GetFloat64(name string) (float64, error)
}
