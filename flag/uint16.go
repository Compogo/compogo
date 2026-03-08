package flag

type Uint16 interface {
	Uint16P(name, shorthand string, value uint16, usage string) *uint16
	Uint16(name string, value uint16, usage string) *uint16
	Uint16VarP(p *uint16, name, shorthand string, value uint16, usage string)
	Uint16Var(p *uint16, name string, value uint16, usage string)
	GetUint16(name string) (uint16, error)
}
