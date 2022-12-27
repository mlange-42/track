package util

import (
	"time"
)

// ParseDate parses a date string
func ParseDate(text string) (time.Time, error) {
	now := time.Now()
	switch text {
	case "today":
		return ToDate(now), nil
	case "tomorrow":
		return ToDate(now).Add(24 * time.Hour), nil
	case "yesterday":
		return ToDate(now).Add(-24 * time.Hour), nil
	}
	return time.ParseInLocation(DateFormat, text, time.Local)
}

// ParseDateTime parses a datetime string. Assumes the local time zone.
func ParseDateTime(text string) (time.Time, error) {
	return time.ParseInLocation(DateTimeFormat, text, time.Local)
}

// ToDate creates a date from a time by setting to 00:00
func ToDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// Date creates a date from a year, a month and a day. Assumes the local time zone.
func Date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

// DateTime creates a datetime
func DateTime(year int, month time.Month, day, hours, minutes, seconds int) time.Time {
	tz := time.Time{}.Location()
	return time.Date(year, month, day, hours, minutes, seconds, 0, tz)
}
