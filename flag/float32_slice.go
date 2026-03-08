package flag

type Float32Slice interface {
	Float32SliceP(name, shorthand string, value []float32, usage string) *[]float32
	Float32Slice(name string, value []float32, usage string) *[]float32
	Float32SliceVarP(p *[]float32, name, shorthand string, value []float32, usage string)
	Float32SliceVar(p *[]float32, name string, value []float32, usage string)
	GetFloat32Slice(name string) ([]float32, error)
}
