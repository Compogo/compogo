package flag

type IntSlice interface {
	IntSliceP(name, shorthand string, value []int, usage string) *[]int
	IntSlice(name string, value []int, usage string) *[]int
	IntSliceVarP(p *[]int, name, shorthand string, value []int, usage string)
	IntSliceVar(p *[]int, name string, value []int, usage string)
	GetIntSlice(name string) ([]int, error)
}
