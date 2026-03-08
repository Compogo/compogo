package flag

type BoolSlice interface {
	BoolSlice(name string, value []bool, usage string) *[]bool
	BoolSliceP(name, shorthand string, value []bool, usage string) *[]bool
	BoolSliceVar(p *[]bool, name string, value []bool, usage string)
	BoolSliceVarP(p *[]bool, name, shorthand string, value []bool, usage string)
	GetBoolSlice(name string) ([]bool, error)
}
