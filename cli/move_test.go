package cli

import (
	"os"
	"testing"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestMove(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	err = track.CreateWorkspace("move-to")
	if err != nil {
		t.Fatal("error creating workspace")
	}

	project := core.NewProject("test", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	record1 := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 4, 5, 0),
		End:     util.DateTime(2001, 2, 3, 5, 5, 0),
		Note:    "Test note with +key=value and +tag and +foo=bar",
	}
	record2 := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 6, 5, 0),
		End:     util.DateTime(2001, 2, 3, 7, 5, 0),
		Note:    "Test note with +tag and +foo=baz",
	}
	err = track.SaveRecord(&record1, false)
	if err != nil {
		t.Fatal("error saving record")
	}
	err = track.SaveRecord(&record2, false)
	if err != nil {
		t.Fatal("error saving record")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"move", "project", "test", "move-to"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	err = track.SwitchWorkspace("move-to")
	if err != nil {
		t.Fatal("error switching workspace")
	}
	assert.Equal(t, "move-to", track.Workspace(), "Should be in new workspace")

	assert.True(t, track.ProjectExists("test"), "Project should exist in new workspace")

	all, err := track.LoadAllRecords()
	if err != nil {
		t.Fatal("error loading records")
	}
	assert.Equal(t, 2, len(all), "Wrong number of records")
}
