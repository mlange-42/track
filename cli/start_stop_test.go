package cli

import (
	"os"
	"testing"

	"github.com/mlange-42/track/core"
	"github.com/stretchr/testify/assert"
)

func TestStartStop(t *testing.T) {
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
	project = core.NewProject("test2", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"start", "test", "Note", "--ago", "60m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"switch", "test2", "Note", "--ago", "30m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"stop"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	all, err := track.LoadAllRecords()
	if err != nil {
		t.Fatal("error loading records")
	}
	assert.Equal(t, 2, len(all), "Wrong number of records")
	assert.Equal(t, all[0].Project, "test", "Wrong record project")
	assert.Equal(t, all[1].Project, "test2", "Wrong record project")
}
