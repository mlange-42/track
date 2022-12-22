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
