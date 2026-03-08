package flag

type Int32 interface {
	Int32P(name, shorthand string, value int32, usage string) *int32
	Int32(name string, value int32, usage string) *int32
	Int32VarP(p *int32, name, shorthand string, value int32, usage string)
	Int32Var(p *int32, name string, value int32, usage string)
	GetInt32(name string) (int32, error)
}
