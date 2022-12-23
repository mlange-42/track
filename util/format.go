package util

import (
	"fmt"
	"time"
)

const (
	// DateFormat for date formatting
	DateFormat = "2006-01-02"
	// FileDateFormat for file name date formatting
	FileDateFormat = "2006-01-02"
	// TimeFormat for time formatting
	TimeFormat = "15:04"
	// FileTimeFormat for file name time formatting
	FileTimeFormat = "15-04"
	// DateTimeFormat for date and time formatting
	DateTimeFormat = "2006-01-02 15:04"
	// JSONTimeFormat for JSON date and time formatting
	JSONTimeFormat = "2006-01-02 15:04:05"
	// NoTime string representation for zero time
	NoTime = " --- "
	// NoDateTime string representation for zero time
	NoDateTime = "      ---       "
)

// FormatDuration formats a duration
func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%.1fhr", d.Hours())
}
