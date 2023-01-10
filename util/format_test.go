package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatDuration(t *testing.T) {
	tt := []struct {
		title    string
		dur      time.Duration
		pad      bool
		expected string
	}{
		{
			title:    "padded hours",
			dur:      time.Hour + 5*time.Minute,
			pad:      true,
			expected: "01:05",
		},
		{
			title:    "padded hours, long numbers",
			dur:      100*time.Hour + 5*time.Minute,
			pad:      true,
			expected: "100:05",
		},
		{
			title:    "no padded hours",
			dur:      time.Hour + 5*time.Minute,
			pad:      false,
			expected: "1:05",
		},
		{
			title:    "no padded hours, long numbers",
			dur:      100*time.Hour + 5*time.Minute,
			pad:      false,
			expected: "100:05",
		},
	}

	for _, test := range tt {
		str := FormatDuration(test.dur, test.pad)
		assert.Equal(t, test.expected, str, "Wrong duration formatting in %s", test.title)
	}
}

func TestFormatTimeWithOffset(t *testing.T) {
	tt := []struct {
		title    string
		time     time.Time
		ref      time.Time
		expected string
	}{
		{
			title:    "same day",
			time:     DateTime(2001, 2, 3, 4, 5, 6),
			ref:      Date(2001, 2, 3),
			expected: "04:05",
		},
		{
			title:    "previous day",
			time:     DateTime(2001, 2, 3, 4, 5, 6),
			ref:      Date(2001, 2, 4),
			expected: "<04:05",
		},
		{
			title:    "next day",
			time:     DateTime(2001, 2, 3, 4, 5, 6),
			ref:      Date(2001, 2, 2),
			expected: "04:05>",
		},
	}

	for _, test := range tt {
		str := FormatTimeWithOffset(test.time, test.ref)
		assert.Equal(t, test.expected, str, "Wrong duration formatting in %s", test.title)
	}
}

func TestFormat(t *testing.T) {
	assert.Equal(
		t,
		"foo baz bar",
		Format("foo {name} bar", map[string]string{"name": "baz"}),
		"Simple replacement not working",
	)
	assert.Equal(
		t,
		"foo baz bar baz",
		Format("foo {name} bar {name}", map[string]string{"name": "baz"}),
		"Repetitions not working",
	)
	assert.Equal(
		t,
		"foo baz bar foo",
		Format("foo {name} bar {name2}", map[string]string{"name": "baz", "name2": "foo"}),
		"Repetitions not working",
	)
}
