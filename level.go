package compogo

//go:generate stringer -type=Level -linecomment

const (
	Panic Level = iota // panic
	Error              // error
	Warn               // warn
	Info               // info
	Debug              // debug
)

type Level uint8
