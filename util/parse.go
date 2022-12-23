package util

import (
	"time"
)

// ParseDate parses a date string
func ParseDate(text string) (time.Time, error) {
	switch text {
	case "today":
		return time.Now(), nil
	case "yesterday":
		return time.Now().Add(-24 * time.Hour), nil
	}
	return time.Parse(DateFormat, text)
}

// ParseDateTime parses a datetime string
func ParseDateTime(text string) (time.Time, error) {
	return time.Parse(DateTimeFormat, text)
}

// Date creates a date
func Date(year int, month time.Month, day int) time.Time {
	tz := time.Time{}.Location()
	return time.Date(year, month, day, 0, 0, 0, 0, tz)
}

// DateTime creates a datetime
func DateTime(year int, month time.Month, day, hours, minutes, seconds int) time.Time {
	tz := time.Time{}.Location()
	return time.Date(year, month, day, hours, minutes, seconds, 0, tz)
}
