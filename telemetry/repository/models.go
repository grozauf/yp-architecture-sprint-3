package repository

import "time"

type TelemetryValue struct {
	Value       float32
	OccuranceAt time.Time
}
