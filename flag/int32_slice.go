package flag

type Int32Slice interface {
	Int32SliceP(name, shorthand string, value []int32, usage string) *[]int32
	Int32Slice(name string, value []int32, usage string) *[]int32
	Int32SliceVarP(p *[]int32, name, shorthand string, value []int32, usage string)
	Int32SliceVar(p *[]int32, name string, value []int32, usage string)
	GetInt32Slice(name string) ([]int32, error)
}
