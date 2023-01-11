package cli

import (
	"os"
	"testing"

	"github.com/mlange-42/track/core"
)

func TestPauseResume(t *testing.T) {
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

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"start", "test", "Note", "--ago", "60m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"pause", "Note", "--ago", "50m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"resume", "--ago", "40m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"pause", "Note", "--ago", "30m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"resume", "--ago", "20m", "--skip"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"pause", "--duration", "10m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

}

func TestResumeLast(t *testing.T) {
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

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"start", "test", "Note", "--ago", "60m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"stop", "--ago", "50m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"resume", "--last", "--skip", "--ago", "40m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"stop", "--ago", "30m"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"resume", "--last"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}

}
