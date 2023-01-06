package core

import (
	"testing"
	"time"

	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestDurationPause(t *testing.T) {
	tt := []struct {
		title       string
		record      Record
		expHasEnded bool
		expIsPaused bool
		expDuration time.Duration
		expPause    time.Duration
	}{
		{
			title: "finished record without pauses",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 9, 0, 0, 0, time.Local),
				Pause:   make([]Pause, 0),
			},
			expHasEnded: true,
			expIsPaused: false,
			expDuration: time.Hour,
			expPause:    0,
		},
		{
			title: "open record without pauses",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     util.NoTime,
				Pause:   make([]Pause, 0),
			},
			expHasEnded: false,
			expIsPaused: false,
			expDuration: time.Hour * 16,
			expPause:    0,
		},
		{
			title: "open record with open pause",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     util.NoTime,
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
						End:   util.NoTime,
					},
				},
			},
			expHasEnded: false,
			expIsPaused: true,
			expDuration: time.Hour * 10,
			expPause:    time.Hour * 6,
		},
	}

	for _, test := range tt {
		end := time.Date(2001, 2, 4, 0, 0, 0, 0, time.Local)
		assert.Equal(t, test.expHasEnded, test.record.HasEnded(), "Wrong HasEnded in %s", test.title)
		assert.Equal(t, test.expIsPaused, test.record.IsPaused(), "Wrong IsPaused in %s", test.title)
		assert.Equal(t, test.expDuration, test.record.Duration(util.NoTime, end), "Wrong duration in %s", test.title)
		assert.Equal(t, test.expPause, test.record.PauseDuration(util.NoTime, end), "Wrong pause duration in %s", test.title)
	}
}

func TestExtractTags(t *testing.T) {
	tt := []struct {
		title   string
		note    string
		expTags []string
	}{
		{
			title:   "no tags",
			note:    "Note without tags",
			expTags: []string{},
		},
		{
			title:   "one tags",
			note:    "Note with a +tag in it",
			expTags: []string{"tag"},
		},
		{
			title:   "two tags",
			note:    "Note with +two +tags in it",
			expTags: []string{"two", "tags"},
		},
		{
			title:   "repeated tags",
			note:    "Note with +two +tags in it +tags +two",
			expTags: []string{"two", "tags"},
		},
	}

	for _, test := range tt {
		tags := ExtractTags(test.note)
		assert.Equal(t, test.expTags, tags, "Failed extracting tags %s", test.title)
	}
}

func TestExtractTagsSlice(t *testing.T) {
	tt := []struct {
		title   string
		note    []string
		expTags []string
	}{
		{
			title:   "no tags",
			note:    []string{"Note without tags"},
			expTags: []string{},
		},
		{
			title:   "one tag",
			note:    []string{"Note with a +tag in it"},
			expTags: []string{"tag"},
		},
		{
			title:   "tag on 2ng line",
			note:    []string{"No tag", "Note with a +tag in it"},
			expTags: []string{"tag"},
		},
		{
			title:   "two tags",
			note:    []string{"Note with +two +tags in it"},
			expTags: []string{"two", "tags"},
		},
		{
			title:   "repeated tags",
			note:    []string{"Note with +two +tags in it +tags +two"},
			expTags: []string{"two", "tags"},
		},
	}

	for _, test := range tt {
		tags := ExtractTagsSlice(test.note)
		assert.Equal(t, test.expTags, tags, "Failed extracting tags %s", test.title)
	}
}
