package flag

type Float64Slice interface {
	Float64SliceP(name, shorthand string, value []float64, usage string) *[]float64
	Float64Slice(name string, value []float64, usage string) *[]float64
	Float64SliceVarP(p *[]float64, name, shorthand string, value []float64, usage string)
	Float64SliceVar(p *[]float64, name string, value []float64, usage string)
	GetFloat64Slice(name string) ([]float64, error)
}
