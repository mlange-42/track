package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationClip(t *testing.T) {
	tt := []struct {
		title    string
		start    time.Time
		end      time.Time
		min      time.Time
		max      time.Time
		expected time.Duration
	}{
		{
			title:    "simple, not clipped",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      DateTime(2001, 2, 3, 5, 5, 0),
			min:      NoTime,
			max:      NoTime,
			expected: time.Hour,
		},
		{
			title:    "simple, clip outside",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      DateTime(2001, 2, 3, 5, 5, 0),
			min:      DateTime(2001, 2, 3, 3, 5, 0),
			max:      DateTime(2001, 2, 3, 6, 5, 0),
			expected: time.Hour,
		},
		{
			title:    "one day, not clipped",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      DateTime(2001, 2, 4, 4, 5, 0),
			min:      NoTime,
			max:      NoTime,
			expected: 24 * time.Hour,
		},
		{
			title:    "one day, start clipped",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      DateTime(2001, 2, 4, 4, 5, 0),
			min:      DateTime(2001, 2, 4, 3, 5, 0),
			max:      NoTime,
			expected: time.Hour,
		},
		{
			title:    "one day, end clipped",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      DateTime(2001, 2, 4, 4, 5, 0),
			min:      NoTime,
			max:      DateTime(2001, 2, 3, 5, 5, 0),
			expected: time.Hour,
		},
		{
			title:    "one day, both clipped",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      DateTime(2001, 2, 4, 4, 5, 0),
			min:      DateTime(2001, 2, 3, 6, 5, 0),
			max:      DateTime(2001, 2, 3, 7, 5, 0),
			expected: time.Hour,
		},
		{
			title:    "one day, clip after",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      DateTime(2001, 2, 4, 4, 5, 0),
			min:      DateTime(2001, 2, 5, 6, 5, 0),
			max:      DateTime(2001, 2, 5, 7, 5, 0),
			expected: 0 * time.Hour,
		},
		{
			title:    "one day, clip before",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      DateTime(2001, 2, 4, 4, 5, 0),
			min:      DateTime(2001, 2, 2, 6, 5, 0),
			max:      DateTime(2001, 2, 2, 7, 5, 0),
			expected: 0 * time.Hour,
		},
		{
			title:    "open end time",
			start:    DateTime(2001, 2, 3, 4, 5, 0),
			end:      NoTime,
			min:      NoTime,
			max:      DateTime(2001, 2, 3, 5, 5, 0),
			expected: time.Hour,
		},
	}

	for _, test := range tt {
		dur := DurationClip(test.start, test.end, test.min, test.max)
		assert.Equal(t, test.expected, dur, "Wrong clipped duration in %s", test.title)
	}
}

func TestMonday(t *testing.T) {
	for i := 1900; i < 2020; i++ {
		date := Date(i, 1, 1)
		monday := Monday(date)
		assert.Equal(t, time.Monday, monday.Weekday(), "Weekday should be monday")
	}
}
