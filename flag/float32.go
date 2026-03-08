package flag

type Float32 interface {
	Float32P(name, shorthand string, value float32, usage string) *float32
	Float32(name string, value float32, usage string) *float32
	Float32VarP(p *float32, name, shorthand string, value float32, usage string)
	Float32Var(p *float32, name string, value float32, usage string)
	GetFloat32(name string) (float32, error)
}
