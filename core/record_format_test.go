package core

import (
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
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
				Tags:    make(map[string]string, 0),
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
				Tags:    make(map[string]string, 0),
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
				Tags: map[string]string{"tag": ""},
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
				t.Fatalf("got unexpected error in %s: %s", test.title, err.Error())
			}
		} else {
			if test.expError {
				t.Fatalf("expected error not raised in %s", test.title)
			}
			assert.Equal(t, test.record, outRecord, "Deserialized record not as expected %s", test.title)
		}
	}
}

func BenchmarkSerialize(b *testing.B) {
	record := fullRecord()
	for i := 0; i < b.N; i++ {
		_ = SerializeRecord(&record, record.Start)
	}
}

func BenchmarkDeserialize(b *testing.B) {
	record := fullRecord()
	text := SerializeRecord(&record, record.Start)

	for i := 0; i < b.N; i++ {
		_, _ = DeserializeRecord(text, record.Start)
	}
}

func BenchmarkRead1k(b *testing.B) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(b, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(b, err, "Error creating Track instance")

	generateDataset(
		&track,
		util.Date(1990, 1, 1),
		2*time.Hour,
		1000,
	)

	for i := 0; i < b.N; i++ {
		_, err = track.LoadAllRecords()
		assert.Nil(b, err, "error loading records")
	}
}

func BenchmarkRead1kAsync(b *testing.B) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(b, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(b, err, "Error creating Track instance")

	generateDataset(
		&track,
		util.Date(1990, 1, 1),
		2*time.Hour,
		1000,
	)

	for i := 0; i < b.N; i++ {
		fn, results, _ := track.AllRecords()
		go fn()
		for res := range results {
			_ = res.Record
		}
	}
}

func BenchmarkRead10k(b *testing.B) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(b, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(b, err, "Error creating Track instance")

	generateDataset(
		&track,
		util.Date(1990, 1, 1),
		2*time.Hour,
		10000,
	)

	for i := 0; i < b.N; i++ {
		_, err = track.LoadAllRecords()
		assert.Nil(b, err, "error loading records")
	}
}

func BenchmarkRead10kAsync(b *testing.B) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(b, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(b, err, "Error creating Track instance")

	generateDataset(
		&track,
		util.Date(1990, 1, 1),
		2*time.Hour,
		10000,
	)

	for i := 0; i < b.N; i++ {
		fn, results, _ := track.AllRecords()
		go fn()
		for res := range results {
			_ = res.Record
		}
	}
}

func fullRecord() Record {
	return Record{
		Project: "test",
		Start:   time.Date(2001, 2, 3, 8, 0, 0, 0, time.Local),
		End:     time.Date(2001, 2, 3, 17, 0, 0, 0, time.Local),
		Pause: []Pause{
			{
				Start: time.Date(2001, 2, 3, 8, 30, 0, 0, time.Local),
				End:   time.Date(2001, 2, 3, 8, 40, 0, 0, time.Local),
				Note:  "Breakfast",
			},
			{
				Start: time.Date(2001, 2, 3, 12, 30, 0, 0, time.Local),
				End:   time.Date(2001, 2, 3, 13, 0, 0, 0, time.Local),
				Note:  "Lunch",
			},
		},
		Note: "Note with a +tag",
		Tags: map[string]string{"tag": ""},
	}
}

func generateDataset(t *Track, start time.Time, step time.Duration, records int) error {
	currTime := start
	for i := 0; i < records; i++ {
		rec := timedRecord(currTime, step/2, rand.Intn(3), rand.Intn(3))
		if err := t.SaveRecord(&rec, false); err != nil {
			return err
		}
		currTime = currTime.Add(step)
	}
	return nil
}

var noteLines = strings.Split(
	`Lorem ipsum dolor sit +amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut +labore et dolore magna aliqua.
A diam sollicitudin +tempor id eu nisl nunc mi ipsum.
Ullamcorper sit amet +risus +nullam eget felis eget.
Cursus euismod quis viverra nibh. Nunc sed blandit libero volutpat sed.
At augue eget arcu dictum varius +duis at consectetur lorem.
Velit euismod in pellentesque +massa placerat duis ultricies.
Risus nec +feugiat in fermentum posuere urna.
Id diam vel quam +elementum +pulvinar +etiam non quam.
Adipiscing bibendum est ultricies +integer quis auctor elit sed.
Sagittis orci a scelerisque purus +semper eget duis.
Elementum tempus egestas sed sed risus pretium.
Sit amet +mattis vulputate enim nulla aliquet porttitor lacus luctus.
Dignissim +sodales ut eu sem integer vitae justo eget.
Posuere urna nec tincidunt +praesent semper feugiat nibh sed pulvinar.
Dignissim cras tincidunt lobortis +feugiat vivamus at augue eget.
Viverra nam libero justo laoreet sit +amet cursus.
Venenatis +lectus +magna +fringilla +urna porttitor rhoncus dolor.
Aliquam id +diam maecenas ultricies.`, "\n")

func timedRecord(start time.Time, duration time.Duration, pauses int, lines int) Record {
	pause := make([]Pause, pauses, pauses)
	note := make([]string, lines, lines)

	if pauses > 0 {
		step := time.Duration(int(duration) / (pauses + 2))
		pauseStart := start.Add(step / 2)
		for i := 0; i < pauses; i++ {
			pause[i] = Pause{
				Start: pauseStart,
				End:   pauseStart.Add(step / 2),
				Note:  "Pause comment",
			}
			pauseStart = pauseStart.Add(step)
		}
	}

	for i := 0; i < lines; i++ {
		note[i] = noteLines[rand.Intn(len(noteLines))]
	}

	return Record{
		Project: "test",
		Start:   start,
		End:     start.Add(duration),
		Pause:   pause,
		Note:    strings.Join(note, "\n"),
		Tags:    map[string]string{},
	}
}
