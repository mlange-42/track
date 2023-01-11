package cli

import (
	"os"
	"testing"
	"time"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/util"
)

func TestStatus(t *testing.T) {
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

	now := time.Now()

	record := core.Record{
		Project: "test",
		Start:   now.Add(-6 * time.Hour),
		End:     now.Add(-5 * time.Hour),
		Note:    "Test note with +key=value and +tag and +foo=bar",
	}
	err = track.SaveRecord(&record, false)
	if err != nil {
		t.Fatal("error saving record")
	}
	record = core.Record{
		Project: "test",
		Start:   now.Add(-4 * time.Hour),
		End:     now.Add(-3 * time.Hour),
		Note:    "Test note with +key=value and +tag and +foo=bar",
	}
	err = track.SaveRecord(&record, false)
	if err != nil {
		t.Fatal("error saving record")
	}
	record = core.Record{
		Project: "test2",
		Start:   now.Add(-2 * time.Hour),
		End:     now.Add(-1 * time.Hour),
		Note:    "Test note with +key=value and +tag and +foo=bar",
	}
	err = track.SaveRecord(&record, false)
	if err != nil {
		t.Fatal("error saving record")
	}
	record = core.Record{
		Project: "test2",
		Start:   now.Add(-30 * time.Minute),
		End:     util.NoTime,
		Note:    "Test note with +key=value and +tag and +foo=bar",
	}
	err = track.SaveRecord(&record, false)
	if err != nil {
		t.Fatal("error saving record")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"status"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"status", "test"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"status", "test2"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}
}
