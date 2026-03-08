package flag

type Int16 interface {
	Int16P(name, shorthand string, value int16, usage string) *int16
	Int16(name string, value int16, usage string) *int16
	Int16VarP(p *int16, name, shorthand string, value int16, usage string)
	Int16Var(p *int16, name string, value int16, usage string)
	GetInt16(name string) (int16, error)
}
