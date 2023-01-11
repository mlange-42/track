package core

import (
	"os"
	"testing"
	"time"

	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestReporter(t *testing.T) {
	dir, err := os.MkdirTemp("", "track-test")
	if err != nil {
		t.Fatal("error creating temporary directory")
	}
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	if err != nil {
		t.Fatal("error creating Track instance")
	}

	project := NewProject("test", "", "T", []string{}, 0, 15)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}
	child := NewProject("child", "test", "T", []string{}, 0, 15)
	err = track.SaveProject(child, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	for i := 0; i < 23; i++ {
		record := Record{
			Project: "child",
			Start:   util.DateTime(2001, 2, 3, i, 0, 0),
			End:     util.DateTime(2001, 2, 3, i, 30, 0),
			Note:    "Test note with +key=value and +tag and +foo=bar",
		}
		err = track.SaveRecord(&record, false)
		if err != nil {
			t.Fatal("error saving record")
		}
	}

	reporter, err := NewReporter(
		&track, []string{}, FilterFunctions{},
		false, util.NoTime, util.NoTime,
	)
	if err != nil {
		t.Fatal("error creating reporter")
	}
	assert.Equal(t, 11*time.Hour+30*time.Minute, reporter.TotalTime["test"], "Wrong total time")

	reporter, err = NewReporter(
		&track, []string{"test", "child"}, FilterFunctions{},
		false, util.NoTime, util.NoTime,
	)
	if err != nil {
		t.Fatal("error creating reporter")
	}
	assert.Equal(t, 11*time.Hour+30*time.Minute, reporter.TotalTime["test"], "Wrong total time")

	reporter, err = NewReporter(
		&track, []string{"foo"}, FilterFunctions{},
		false, util.NoTime, util.NoTime,
	)
	assert.NotNil(t, err, "expecting error on invalid project")
}
