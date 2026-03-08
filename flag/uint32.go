package flag

type Uint32 interface {
	Uint32P(name, shorthand string, value uint32, usage string) *uint32
	Uint32(name string, value uint32, usage string) *uint32
	Uint32VarP(p *uint32, name, shorthand string, value uint32, usage string)
	Uint32Var(p *uint32, name string, value uint32, usage string)
	GetUint32(name string) (uint32, error)
}
