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
	return time.Date(year, month, day, hours, minutes, seconds, 0, time.Local)
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

	start, err = ParseTimeWithOffset(parts[0], date)
	if err != nil {
		return
	}

	if parts[1] != "?" {
		end, err = ParseTimeWithOffset(parts[1], date)
		if err != nil {
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

// ParseTimeWithOffset parses a time with offset markers
func ParseTimeWithOffset(text string, date time.Time) (time.Time, error) {
	dayOffset := 0
	if strings.HasPrefix(text, PrevDayPrefix) {
		text = strings.TrimPrefix(text, PrevDayPrefix)
		dayOffset--
	}
	if strings.HasSuffix(text, NextDaySuffix) {
		text = strings.TrimSuffix(text, NextDaySuffix)
		dayOffset++
	}
	t, err := time.ParseInLocation(TimeFormat, text, time.Local)
	if err != nil {
		return NoTime, err
	}
	t = DateAndTime(date, t)
	if dayOffset != 0 {
		t = t.Add(time.Duration(int64(dayOffset*24) * int64(time.Hour)))
	}
	return t, nil
}

// DateAndTime combines a date with a time
func DateAndTime(d, t time.Time) time.Time {
	return time.Date(
		d.Year(), d.Month(), d.Day(),
		t.Hour(), t.Minute(), t.Second(), 0, time.Local,
	)
}
