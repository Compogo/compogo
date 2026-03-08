package logger

import "github.com/Compogo/compogo/types"

//go:generate stringer -type=Level

const (
	Panic Level = iota
	Error
	Warn
	Info
	Debug
)

type Level uint8

var Levels = types.NewMapper[Level](Panic, Error, Warn, Info, Debug)
