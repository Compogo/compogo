package logger

import "github.com/Compogo/types/mapper"

//go:generate stringer -type=Level

const (
	Panic Level = iota
	Error
	Warn
	Info
	Debug
)

type Level uint8

var Levels = mapper.NewMapper[Level](Panic, Error, Warn, Info, Debug)
