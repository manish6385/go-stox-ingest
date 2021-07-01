package models

import "time"

type Order struct {
	ID       int
	Customer string
	Date     time.Time
}
