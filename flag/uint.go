package flag

type Uint interface {
	UintP(name, shorthand string, value uint, usage string) *uint
	Uint(name string, value uint, usage string) *uint
	UintVarP(p *uint, name, shorthand string, value uint, usage string)
	UintVar(p *uint, name string, value uint, usage string)
}
