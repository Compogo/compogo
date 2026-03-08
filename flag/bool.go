package flag

type Bool interface {
	Bool(name string, value bool, usage string) *bool
	BoolP(name, shorthand string, value bool, usage string) *bool
	BoolVar(p *bool, name string, value bool, usage string)
	BoolVarP(p *bool, name, shorthand string, value bool, usage string)
	GetBool(name string) (bool, error)
}
