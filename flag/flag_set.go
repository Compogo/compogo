package flag

type FlagSet interface {
	Bool
	BoolSlice
	Bytes
	Count
	Duration
	DurationSlice
	Float32
	Float32Slice
	Float64
	Float64Slice
	Int
	Int8
	Int16
	Int32
	Int32Slice
	Int64
	Int64Slice
	IntSlice
	IP
	IPSlice
	IPMask
	IPNet
	String
	StringArray
	StringSlice
	StringToInt
	StringToInt64
	StringToString
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	UintSlice
}
