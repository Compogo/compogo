package flag

type Bytes interface {
	BytesBase64(name string, value []byte, usage string) *[]byte
	BytesBase64P(name, shorthand string, value []byte, usage string) *[]byte
	BytesBase64VarP(p *[]byte, name, shorthand string, value []byte, usage string)
	BytesBase64Var(p *[]byte, name string, value []byte, usage string)
	GetBytesBase64(name string) ([]byte, error)

	BytesHexP(name, shorthand string, value []byte, usage string) *[]byte
	BytesHex(name string, value []byte, usage string) *[]byte
	BytesHexVarP(p *[]byte, name, shorthand string, value []byte, usage string)
	BytesHexVar(p *[]byte, name string, value []byte, usage string)
	GetBytesHex(name string) ([]byte, error)
}
