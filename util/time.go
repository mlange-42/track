package util

import "time"

// NoTime is a zero time
var NoTime time.Time = time.Time{}

// DurationClip calculated the duration a time span, clipped by another time span
func DurationClip(start, end, min, max time.Time) time.Duration {
	if end.IsZero() {
		end = time.Now()
	}
	if !min.IsZero() {
		if end.Before(min) {
			return 0 * time.Second
		}
		if start.Before(min) {
			start = min
		}
	}
	if !max.IsZero() {
		if start.After(max) {
			return 0 * time.Second
		}
		if end.After(max) {
			end = max
		}
	}

	return end.Sub(start)
}

// Monday returns the monday of the week of the given date
func Monday(date time.Time) time.Time {
	weekDay := (int(date.Weekday()) + 6) % 7
	return date.Add(time.Duration(-weekDay * 24 * int(time.Hour)))
}
