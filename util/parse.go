package util

import (
	"fmt"
	"strings"
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

// ParseTimeRange parses a time range string. Assumes the local time zone.
func ParseTimeRange(text string, date time.Time) (start, end time.Time, err error) {
	parts := strings.Split(text, "-")
	if len(parts) != 2 {
		err = fmt.Errorf("invalid time range: must be 2 parts")
		return
	}
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}

	start, err = time.ParseInLocation(TimeFormat, parts[0], time.Local)
	if err != nil {
		return
	}
	start = DateAndTime(date, start)
	nextDay := false
	if strings.HasPrefix(parts[1], "+") {
		parts[1] = strings.TrimPrefix(parts[1], "+")
		nextDay = true
	}

	if parts[1] != "?" {
		end, err = time.ParseInLocation(TimeFormat, parts[1], time.Local)
		if err == nil {
			end = DateAndTime(date, end)
			if nextDay {
				end = end.Add(24 * time.Hour)
			}
		} else {
			var dur time.Duration
			dur, err = time.ParseDuration(parts[1])
			if err != nil {
				return
			}
			end = start.Add(dur)
		}

		if start.After(end) {
			err = fmt.Errorf("invalid time range: start must be before end")
			return
		}
	}

	return start, end, nil
}

// DateAndTime combines a date with a time
func DateAndTime(d, t time.Time) time.Time {
	return time.Date(
		d.Year(), d.Month(), d.Day(),
		t.Hour(), t.Minute(), t.Second(), 0, time.Local,
	)
}
