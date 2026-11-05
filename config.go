package compogo

import "time"

type Config struct {
	Name      string
	PID       uint64
	StartTime time.Time
}
