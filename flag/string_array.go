package flag

type StringArray interface {
	StringArrayP(name, shorthand string, value []string, usage string) *[]string
	StringArray(name string, value []string, usage string) *[]string
	StringArrayVarP(p *[]string, name, shorthand string, value []string, usage string)
	StringArrayVar(p *[]string, name string, value []string, usage string)
	GetStringArray(name string) ([]string, error)
}
