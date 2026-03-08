package flag

type String interface {
	String(name string, value string, usage string) *string
	StringP(name, shorthand string, value string, usage string) *string
	StringVar(p *string, name string, value string, usage string)
	StringVarP(p *string, name, shorthand string, value string, usage string)
}
