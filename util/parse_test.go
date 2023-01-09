package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTimeRange(t *testing.T) {
	date := Date(2022, 12, 29)

	tt := []struct {
		title    string
		text     string
		expStart time.Time
		expEnd   time.Time
		expErr   bool
	}{
		{
			title:    "Empty string",
			text:     "",
			expStart: NoTime,
			expEnd:   NoTime,
			expErr:   true,
		},
		{
			title:    "No end time",
			text:     "12:30 - ?",
			expStart: date.Add(time.Duration(time.Minute * (12*60 + 30))),
			expEnd:   NoTime,
			expErr:   false,
		},
		{
			title:    "Normal end time",
			text:     "12:30 - 13:15",
			expStart: date.Add(time.Duration(time.Minute * (12*60 + 30))),
			expEnd:   date.Add(time.Duration(time.Minute * (13*60 + 15))),
			expErr:   false,
		},
		{
			title:    "No spaces",
			text:     "12:30-13:15",
			expStart: date.Add(time.Duration(time.Minute * (12*60 + 30))),
			expEnd:   date.Add(time.Duration(time.Minute * (13*60 + 15))),
			expErr:   false,
		},
		{
			title:    "End time next day",
			text:     "12:30 - 0:30>",
			expStart: date.Add(time.Duration(time.Minute * (12*60 + 30))),
			expEnd:   date.Add(time.Duration(time.Minute * (24*60 + 30))),
			expErr:   false,
		},
		{
			title:    "Start time previous day",
			text:     "<23:30 - 2:30",
			expStart: date.Add(time.Duration(time.Minute * (-30))),
			expEnd:   date.Add(time.Duration(time.Minute * (2*60 + 30))),
			expErr:   false,
		},
		{
			title:    "End time by duration",
			text:     "12:30 - 45m",
			expStart: date.Add(time.Duration(time.Minute * (12*60 + 30))),
			expEnd:   date.Add(time.Duration(time.Minute * (13*60 + 15))),
			expErr:   false,
		},
		{
			title:    "End before start",
			text:     "12:30 - 10:30",
			expStart: date.Add(time.Duration(time.Minute * (12*60 + 30))),
			expEnd:   date.Add(time.Duration(time.Minute * (10*60 + 30))),
			expErr:   true,
		},
	}

	for _, test := range tt {
		start, end, err := ParseTimeRange(test.text, date)
		if start != test.expStart {
			t.Fatalf(
				"%s: Mismatch start. Got %s, exp. %s",
				test.title,
				start.Format(DateTimeFormat),
				test.expStart.Format(DateTimeFormat),
			)
		}
		if end != test.expEnd {
			t.Fatalf(
				"%s: Mismatch end. Got %s, exp. %s",
				test.title,
				end.Format(DateTimeFormat),
				test.expEnd.Format(DateTimeFormat),
			)
		}
		if (err != nil) != test.expErr {
			t.Fatalf(
				"%s: Mismatch end. Got %t, exp. %t (%v)",
				test.title,
				(err != nil),
				test.expErr,
				err,
			)
		}
	}
}

func TestParseDate(t *testing.T) {
	today := ToDate(time.Now())

	tt := []struct {
		title   string
		text    string
		expDate time.Time
	}{
		{
			title:   "today",
			text:    "today",
			expDate: today,
		},
		{
			title:   "yesterday",
			text:    "yesterday",
			expDate: today.Add(-24 * time.Hour),
		},
		{
			title:   "tomorrow",
			text:    "tomorrow",
			expDate: today.Add(24 * time.Hour),
		},
		{
			title:   "date",
			text:    "2022-12-31",
			expDate: Date(2022, 12, 31),
		},
	}

	for _, test := range tt {
		date, err := ParseDate(test.text)
		assert.Nil(t, err, "Error parsing date in %s", test.title)
		assert.Equal(t, test.expDate, date, "Wrong date in %s", test.title)
	}
}

func BenchmarkParseTimeRange(b *testing.B) {
	today := ToDate(time.Now())
	text := "10:00 - 18:00"

	for i := 0; i < b.N; i++ {
		_, _, _ = ParseTimeRange(text, today)
	}
}

func BenchmarkParseTimeRangeOffset(b *testing.B) {
	today := ToDate(time.Now())
	text := "<10:00 - 18:00"

	for i := 0; i < b.N; i++ {
		_, _, _ = ParseTimeRange(text, today)
	}
}

func BenchmarkParseTimeOffset(b *testing.B) {
	today := ToDate(time.Now())
	text := "00:31>"

	for i := 0; i < b.N; i++ {
		_, _ = ParseTimeWithOffset(text, today)
	}
}
