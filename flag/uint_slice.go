package flag

type UintSlice interface {
	UintSliceP(name, shorthand string, value []uint, usage string) *[]uint
	UintSlice(name string, value []uint, usage string) *[]uint
	UintSliceVarP(p *[]uint, name, shorthand string, value []uint, usage string)
	UintSliceVar(p *[]uint, name string, value []uint, usage string)
	GetUintSlice(name string) ([]uint, error)
}
