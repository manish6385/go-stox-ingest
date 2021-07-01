package models

import (
	"time"
)

type Holiday struct {
	ID      int
	Holiday string
	Date    time.Time
	Weekday time.Weekday
}
