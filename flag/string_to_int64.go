package flag

type StringToInt64 interface {
	StringToInt64P(name, shorthand string, value map[string]int64, usage string) *map[string]int64
	StringToInt64(name string, value map[string]int64, usage string) *map[string]int64
	StringToInt64VarP(p *map[string]int64, name, shorthand string, value map[string]int64, usage string)
	StringToInt64Var(p *map[string]int64, name string, value map[string]int64, usage string)
	GetStringToInt64(name string) (map[string]int64, error)
}
