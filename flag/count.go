package flag

type Count interface {
	CountP(name, shorthand string, usage string) *int
	Count(name string, usage string) *int
	CountVarP(p *int, name, shorthand string, usage string)
	CountVar(p *int, name string, usage string)
	GetCount(name string) (int, error)
}
