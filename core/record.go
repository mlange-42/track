package core

import "time"

// Record holds and manipulates data for a record
type Record struct {
	Project Project
	Note    string
	Start   time.Time
	End     time.Time
}
