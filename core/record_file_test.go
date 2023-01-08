package core

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/mlange-42/track/fs"
	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestSaveLoadRecord(t *testing.T) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(t, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(t, err, "Error creating Track instance")

	record1 := Record{
		Project: "test",
		Start:   time.Date(2001, 2, 3, 4, 5, 0, 0, time.Local),
		End:     time.Date(2001, 2, 3, 4, 15, 0, 0, time.Local),
		Note:    "Note with +tag",
		Tags:    []string{"tag"},
		Pause: []Pause{
			{
				Start: time.Date(2001, 2, 3, 4, 8, 0, 0, time.Local),
				End:   time.Date(2001, 2, 3, 4, 9, 0, 0, time.Local),
				Note:  "Pause note",
			},
		},
	}
	record2 := Record{
		Project: "test2",
		Start:   time.Date(2001, 6, 2, 23, 0, 0, 0, time.Local),
		End:     time.Date(2001, 6, 3, 4, 0, 0, 0, time.Local),
		Note:    "Note with +tag",
		Tags:    []string{"tag"},
		Pause: []Pause{
			{
				Start: time.Date(2001, 6, 3, 1, 0, 0, 0, time.Local),
				End:   time.Date(2001, 6, 3, 2, 0, 0, 0, time.Local),
				Note:  "Pause note",
			},
		},
	}
	record3 := Record{
		Project: "test",
		Start:   time.Date(2001, 6, 3, 11, 0, 0, 0, time.Local),
		End:     util.NoTime,
		Note:    "Note with +tag",
		Tags:    []string{"tag"},
		Pause: []Pause{
			{
				Start: time.Date(2001, 6, 3, 12, 0, 0, 0, time.Local),
				End:   time.Date(2001, 6, 3, 13, 0, 0, 0, time.Local),
				Note:  "Pause note",
			},
		},
	}

	err = track.SaveRecord(&record1, false)
	assert.Nil(t, err, "Error saving record")
	err = track.SaveRecord(&record2, false)
	assert.Nil(t, err, "Error saving record")
	err = track.SaveRecord(&record3, false)
	assert.Nil(t, err, "Error saving record")

	assert.True(t, fs.FileExists(track.RecordPath(record1.Start)), "File must exist")
	assert.True(t, fs.FileExists(track.RecordPath(record2.Start)), "File must exist")
	assert.True(t, fs.FileExists(track.RecordPath(record3.Start)), "File must exist")

	newRecord, err := track.LoadRecord(record1.Start)
	assert.Nil(t, err, "Error loading record")
	assert.Equal(t, record1, newRecord, "Loaded record not equal to saved record")

	latestRecord, err := track.LatestRecord()
	assert.Nil(t, err, "Error loading record")
	assert.Equal(t, record3, *latestRecord, "Loaded record not equal to saved record")

	openRecord, err := track.OpenRecord()
	assert.Nil(t, err, "Error loading record")
	assert.Equal(t, record3, *openRecord, "Loaded record not equal to saved record")

	for i := 0; i < 25; i++ {
		allRecords, err := track.LoadAllRecords()
		assert.Nil(t, err, "Error loading all records")
		assert.Equal(t, []Record{record1, record2, record3}, allRecords, "Loaded record not equal to saved record")
	}

	err = track.DeleteRecord(&record1)
	assert.Nil(t, err, "Error deleting record")
	assert.False(t, fs.FileExists(track.RecordPath(record1.Start)), "File must exist")
}

func TestStartStopRecord(t *testing.T) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(t, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(t, err, "Error creating Track instance")

	start := time.Now().Round(time.Minute).Add(-time.Hour)
	record, err := track.StartRecord("test", "", []string{}, start)
	assert.Nil(t, err, "Error starting record")

	openRecord, err := track.OpenRecord()
	assert.Nil(t, err, "Error loading record")
	assert.Equal(t, record, *openRecord, "Loaded record not equal to saved record")

	stopped, err := track.StopRecord(start.Add(time.Hour))
	assert.Nil(t, err, "Error loading record")

	openRecord, err = track.OpenRecord()
	assert.Nil(t, err, "Error loading record")
	assert.Nil(t, openRecord, "Loaded record not equal to saved record")

	lastRecord, err := track.LatestRecord()
	assert.Nil(t, err, "Error loading record")
	assert.Equal(t, stopped, lastRecord, "Loaded record not equal to saved record")
}
