package cli

import (
	"os"
	"testing"

	"github.com/mlange-42/track/core"
	"github.com/mlange-42/track/util"
	"github.com/stretchr/testify/assert"
)

func TestEditConfig(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	util.SkipEditingForTests = true

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"edit", "config"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("error executing command: %s", err.Error())
	}
}

func TestEditProject(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	util.SkipEditingForTests = true

	project := core.NewProject("test", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"edit", "project", "test"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("error executing command: %s", err.Error())
	}
}

func TestEditProjectRenameArchive(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	util.SkipEditingForTests = true

	project := core.NewProject("test", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}
	childProject := core.NewProject("child", "test", "t", []string{}, 15, 0)
	err = track.SaveProject(childProject, false)
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

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"edit", "project", "test", "--rename", "other"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("error executing command: %s", err.Error())
	}
	assert.False(t, track.ProjectExists("test"), "Project should not exist")
	assert.True(t, track.ProjectExists("other"), "Project should exist")

	newRec, err := track.LoadRecord(record.Start)
	if err != nil {
		t.Fatal("error loading record")
	}
	assert.Equal(t, "other", newRec.Project, "record should be in renamed project")

	newProj, err := track.LoadProjectByName("child")
	if err != nil {
		t.Fatal("error loading project")
	}
	assert.Equal(t, "other", newProj.Parent, "parent should be the renamed project")

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"edit", "project", "other", "--archive"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("error executing command: %s", err.Error())
	}
	newProj, err = track.LoadProjectByName("other")
	if err != nil {
		t.Fatal("error loading project")
	}
	assert.True(t, newProj.Archived, "project should be archived")
}

func TestEditRecordDay(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	util.SkipEditingForTests = true

	project := core.NewProject("test", "", "t", []string{}, 15, 0)
	err = track.SaveProject(project, false)
	if err != nil {
		t.Fatal("error saving project")
	}

	record1 := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 4, 5, 0),
		End:     util.DateTime(2001, 2, 3, 5, 5, 0),
	}
	err = track.SaveRecord(&record1, false)
	if err != nil {
		t.Fatal("error saving record")
	}
	record2 := core.Record{
		Project: "test",
		Start:   util.DateTime(2001, 2, 3, 8, 5, 0),
		End:     util.DateTime(2001, 2, 3, 13, 5, 0),
	}
	err = track.SaveRecord(&record2, false)
	if err != nil {
		t.Fatal("error saving record")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"edit", "record", "2001-02-03", "04:05"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("error executing command: %s", err.Error())
	}

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"edit", "day", "2001-02-03"})

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("error executing command: %s", err.Error())
	}
}
