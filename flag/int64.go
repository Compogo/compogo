package flag

type Int64 interface {
	Int64P(name, shorthand string, value int64, usage string) *int64
	Int64(name string, value int64, usage string) *int64
	Int64VarP(p *int64, name, shorthand string, value int64, usage string)
	Int64Var(p *int64, name string, value int64, usage string)
	GetInt64(name string) (int64, error)
}
