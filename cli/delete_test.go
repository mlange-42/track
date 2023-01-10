package cli

import (
	"os"
	"strings"
	"testing"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/out"
	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	project := core.NewProject("test", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	record := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 4, 5, 0),
		End:     util.DateTime(2001, 2, 3, 5, 5, 0),
	}
	err = track.SaveRecord(&record, false)
	if err != nil {
		t.Fatal("error saving record")
	}

	_, err = track.LoadRecord(record.Start)
	if err != nil {
		t.Fatal("error loading record - record should exist")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"delete", "record", "2001-02-03", "04:05"})

	out.StdIn = strings.NewReader("y")
	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	_, err = track.LoadRecord(record.Start)
	assert.NotNil(t, err, "expecting error on loading deleted record")

	err = track.SaveRecord(&record, false)
	if err != nil {
		t.Fatal("error saving record")
	}

	_, err = track.LoadRecord(record.Start)
	if err != nil {
		t.Fatal("error loading record - record should exist")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"delete", "project", "test"})

	out.StdIn = strings.NewReader("yes!")
	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	_, err = track.LoadRecord(record.Start)
	assert.NotNil(t, err, "expecting error on loading deleted record")

	assert.False(t, track.ProjectExists("test"), "project should not exist")
}
