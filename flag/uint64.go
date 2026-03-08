package flag

type Uint64 interface {
	Uint64P(name, shorthand string, value uint64, usage string) *uint64
	Uint64(name string, value uint64, usage string) *uint64
	Uint64VarP(p *uint64, name, shorthand string, value uint64, usage string)
	Uint64Var(p *uint64, name string, value uint64, usage string)
	GetUint64(name string) (uint64, error)
}
