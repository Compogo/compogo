package flag

type StringSlice interface {
	StringSliceP(name, shorthand string, value []string, usage string) *[]string
	StringSlice(name string, value []string, usage string) *[]string
	StringSliceVarP(p *[]string, name, shorthand string, value []string, usage string)
	StringSliceVar(p *[]string, name string, value []string, usage string)
	GetStringSlice(name string) ([]string, error)
}
