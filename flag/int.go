package flag

type Int interface {
	IntP(name, shorthand string, value int, usage string) *int
	Int(name string, value int, usage string) *int
	IntVarP(p *int, name, shorthand string, value int, usage string)
	IntVar(p *int, name string, value int, usage string)
	GetInt(name string) (int, error)
}
