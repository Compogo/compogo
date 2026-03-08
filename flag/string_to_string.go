package flag

type StringToString interface {
	StringToStringP(name, shorthand string, value map[string]string, usage string) *map[string]string
	StringToString(name string, value map[string]string, usage string) *map[string]string
	StringToStringVarP(p *map[string]string, name, shorthand string, value map[string]string, usage string)
	StringToStringVar(p *map[string]string, name string, value map[string]string, usage string)
	GetStringToString(name string) (map[string]string, error)
}
