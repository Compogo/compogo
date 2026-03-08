package flag

type StringToInt interface {
	StringToIntP(name, shorthand string, value map[string]int, usage string) *map[string]int
	StringToInt(name string, value map[string]int, usage string) *map[string]int
	StringToIntVarP(p *map[string]int, name, shorthand string, value map[string]int, usage string)
	StringToIntVar(p *map[string]int, name string, value map[string]int, usage string)
	GetStringToInt(name string) (map[string]int, error)
}
