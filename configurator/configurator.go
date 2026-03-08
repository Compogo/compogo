package configurator

import "time"

type Configurator interface {
	GetString(string) string
	GetBool(string) bool
	GetInt(string) int
	GetInt8(string) int8
	GetInt16(string) int16
	GetInt32(string) int32
	GetInt64(string) int64
	GetUint(string) uint
	GetUint8(string) uint8
	GetUint16(string) uint16
	GetUint32(string) uint32
	GetUint64(string) uint64
	GetFloat32(string) float32
	GetFloat64(string) float64
	GetTime(string) time.Time
	GetDuration(string) time.Duration
	GetIntSlice(string) []int
	GetStringSlice(string) []string
	GetStringMap(string) map[string]interface{}
	GetStringMapString(string) map[string]string
	GetStringMapStringSlice(string) map[string][]string
	GetSizeInBytes(string) uint
	SetDefault(string, interface{})
	ReadConfig() error
}
