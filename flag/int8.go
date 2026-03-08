package flag

type Int8 interface {
	Int8P(name, shorthand string, value int8, usage string) *int8
	Int8(name string, value int8, usage string) *int8
	Int8VarP(p *int8, name, shorthand string, value int8, usage string)
	Int8Var(p *int8, name string, value int8, usage string)
	GetInt8(name string) (int8, error)
}
