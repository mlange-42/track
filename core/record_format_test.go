package core

import (
	"testing"
	"time"

	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestSerializeDeserialize(t *testing.T) {
	tt := []struct {
		title    string
		time     time.Time
		record   Record
		text     string
		expError bool
	}{
		{
			title: "minimal record",
			time:  util.Date(2001, 2, 3),
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 9, 0, 0, 0, time.Local),
				Pause:   make([]Pause, 0),
			},
			text: `08:00 - 09:00
    test
`,
			expError: false,
		},
		{
			title: "minimal record with open end time",
			time:  util.Date(2001, 2, 3),
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     util.NoTime,
				Pause:   make([]Pause, 0),
			},
			text: `08:00 - ?
    test
`,
			expError: false,
		},
		{
			title: "record with pauses, note and tags",
			time:  util.Date(2001, 2, 3),
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     util.NoTime,
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 8, 30, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 8, 40, 0, 0, time.Local),
						Note:  "Breakfast",
					},
					{
						Start: time.Date(2001, 2, 3, 12, 30, 0, 0, time.Local),
						End:   util.NoTime,
						Note:  "Lunch",
					},
				},
				Note: "Note with a +tag",
				Tags: []string{"tag"},
			},
			text: `08:00 - ?
    - 08:30 - 10m0s / Breakfast
    - 12:30 - ? / Lunch
    test

Note with a +tag
`,
			expError: false,
		},
	}

	for _, test := range tt {
		outText := SerializeRecord(&test.record, test.time)
		assert.Equal(t, test.text, outText, "Serialized string not as expected %s", test.title)

		outRecord, err := DeserializeRecord(test.text, test.time)
		if err != nil {
			if !test.expError {
				t.Fatalf("got unexpected error in %s", test.title)
			}
		} else {
			if test.expError {
				t.Fatalf("expected error not raised in %s", test.title)
			}
			assert.Equal(t, test.record, outRecord, "Deserialized record not as expected %s", test.title)
		}
	}
}
