package flag

type Uint8 interface {
	Uint8P(name, shorthand string, value uint8, usage string) *uint8
	Uint8(name string, value uint8, usage string) *uint8
	Uint8VarP(p *uint8, name, shorthand string, value uint8, usage string)
	Uint8Var(p *uint8, name string, value uint8, usage string)
	GetUint8(name string) (uint8, error)
}
