package core

import (
	"testing"
	"time"

	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestCheckRecord(t *testing.T) {
	projectTags := NewProject("test", "", "t", []string{"a"}, 15, 0)
	projectNoTags := NewProject("test", "", "t", []string{}, 15, 0)

	tt := []struct {
		title    string
		record   Record
		project  Project
		expError bool
	}{
		{
			title: "fully valid record",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 9, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 10, 0, 0, 0, time.Local),
					},
					{
						Start: time.Date(2001, 2, 3, 12, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 13, 0, 0, 0, time.Local),
					},
				},
			},
			project:  projectNoTags,
			expError: false,
		},
		{
			title: "fully valid record with tags",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Tags:    map[string]string{"a": "value"},
				Pause:   []Pause{},
			},
			project:  projectTags,
			expError: false,
		},
		{
			title: "record with missing tag value",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Tags:    map[string]string{"a": ""},
				Pause:   []Pause{},
			},
			project:  projectTags,
			expError: true,
		},
		{
			title: "record with missing tag",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Tags:    map[string]string{"b": "value"},
				Pause:   []Pause{},
			},
			project:  projectTags,
			expError: true,
		},
		{
			title: "ends before start",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 7, 0, 0, 0, time.Local),
				Pause:   []Pause{},
			},
			expError: true,
		},
		{
			title: "pause too early",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 7, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 10, 0, 0, 0, time.Local),
					},
				},
			},
			expError: true,
		},
		{
			title: "pause too late",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 9, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 19, 0, 0, 0, time.Local),
					},
				},
			},
			expError: true,
		},
		{
			title: "pause ends before start",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 11, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 10, 0, 0, 0, time.Local),
					},
				},
			},
			expError: true,
		},
		{
			title: "pause open but record not",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 11, 0, 0, 0, time.Local),
						End:   util.NoTime,
					},
				},
			},
			expError: true,
		},
		{
			title: "pauses overlap",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 10, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 12, 0, 0, 0, time.Local),
					},
					{
						Start: time.Date(2001, 2, 3, 11, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 13, 0, 0, 0, time.Local),
					},
				},
			},
			expError: true,
		},
		{
			title: "pauses not chronologically",
			record: Record{
				Project: "test",
				Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
				End:     time.Date(2001, 2, 3, 18, 0, 0, 0, time.Local),
				Pause: []Pause{
					{
						Start: time.Date(2001, 2, 3, 12, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 13, 0, 0, 0, time.Local),
					},
					{
						Start: time.Date(2001, 2, 3, 10, 0, 0, 0, time.Local),
						End:   time.Date(2001, 2, 3, 11, 0, 0, 0, time.Local),
					},
				},
			},
			expError: true,
		},
	}

	for _, test := range tt {
		err := test.record.Check(&test.project)
		if err != nil {
			if !test.expError {
				t.Fatalf("got unexpected error in  %s: %s", test.title, err.Error())
			}
		} else {
			if test.expError {
				t.Fatalf("expected error not raised in %s", test.title)
			}
		}
	}
}

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
		expTags map[string]string
	}{
		{
			title:   "no tags",
			note:    "Note without tags",
			expTags: map[string]string{},
		},
		{
			title:   "one tags",
			note:    "Note with a +tag in it",
			expTags: map[string]string{"tag": ""},
		},
		{
			title:   "two tags",
			note:    "Note with +two +tags in it",
			expTags: map[string]string{"two": "", "tags": ""},
		},
		{
			title:   "repeated tags",
			note:    "Note with +two +tags in it +tags +two",
			expTags: map[string]string{"two": "", "tags": ""},
		},
		{
			title:   "tags with values",
			note:    "Note with +two=foo +tags= in it",
			expTags: map[string]string{"two": "foo", "tags": ""},
		},
	}

	for _, test := range tt {
		tags, err := ExtractTags(test.note)
		assert.Nil(t, err, "Error extracting tags")
		assert.Equal(t, test.expTags, tags, "Failed extracting tags %s", test.title)
	}
}

func TestExtractTagsSlice(t *testing.T) {
	tt := []struct {
		title   string
		note    []string
		expTags map[string]string
	}{
		{
			title:   "no tags",
			note:    []string{"Note without tags"},
			expTags: map[string]string{},
		},
		{
			title:   "one tag",
			note:    []string{"Note with a +tag in it"},
			expTags: map[string]string{"tag": ""},
		},
		{
			title:   "tag on 2ng line",
			note:    []string{"No tag", "Note with a +tag in it"},
			expTags: map[string]string{"tag": ""},
		},
		{
			title:   "two tags",
			note:    []string{"Note with +two +tags in it"},
			expTags: map[string]string{"two": "", "tags": ""},
		},
		{
			title:   "repeated tags",
			note:    []string{"Note with +two +tags in it +tags +two"},
			expTags: map[string]string{"two": "", "tags": ""},
		},
		{
			title:   "tags with value",
			note:    []string{"Note with +two=foo +tags= in it"},
			expTags: map[string]string{"two": "foo", "tags": ""},
		},
	}

	for _, test := range tt {
		tags, err := ExtractTagsSlice(test.note)
		assert.Nil(t, err, "Error extracting tags")
		assert.Equal(t, test.expTags, tags, "Failed extracting tags %s", test.title)
	}
}

func BenchmarkExtractTags(b *testing.B) {
	text := "a test text with a +tag and a +key=value pair"

	for i := 0; i < b.N; i++ {
		_, _ = ExtractTags(text)
	}
}

func BenchmarkExtractTagsSlice(b *testing.B) {
	text := []string{"a test text with a +tag and a +key=value pair"}

	for i := 0; i < b.N; i++ {
		_, _ = ExtractTagsSlice(text)
	}
}
