package compogo

import (
	"github.com/Compogo/types/mapper"
	"github.com/Compogo/types/set"
)

//go:generate stringer -type=Level

const (
	Panic Level = iota
	Error
	Warn
	Info
	Debug
)

type Level uint8

var AllLevels = mapper.NewMapper[Level](Panic, Error, Warn, Info, Debug)

var AllErrorLevels = set.NewSet[Level](Panic, Error, Warn)

var AllInfoLevels = set.NewSet[Level](Info, Debug)
