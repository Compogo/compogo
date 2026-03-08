package flag

type Int64Slice interface {
	Int64SliceP(name, shorthand string, value []int64, usage string) *[]int64
	Int64Slice(name string, value []int64, usage string) *[]int64
	Int64SliceVarP(p *[]int64, name, shorthand string, value []int64, usage string)
	Int64SliceVar(p *[]int64, name string, value []int64, usage string)
	GetInt64Slice(name string) ([]int64, error)
}
