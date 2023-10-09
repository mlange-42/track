package cli

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mlange-42/track/core"
	"github.com/stretchr/testify/assert"
)

func TestCreateWorkspace(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"create", "workspace", "test-ws"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}
	assert.True(t, track.WorkspaceExists("test-ws"), "Workspace should exist")

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"create", "workspace", "test-ws"})

	err = cmd.Execute()
	assert.NotNil(t, err, "should fail with workspace already exists error")
}

func TestCreateProjectSimple(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"create", "project", "test"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}
	assert.True(t, track.ProjectExists("test"), "Project should exist")

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{"create", "project", "test"})

	err = cmd.Execute()
	assert.NotNil(t, err, "should fail with project already exists error")
}

func TestCreateProjectFlags(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"create", "project", "test"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}
	assert.True(t, track.ProjectExists("test"), "Project should exist")

	ref := core.NewProject("child", "test", "C", []string{"tag1", "tag2"}, 15, 0)

	cmd = RootCommand(track, "")
	cmd.SetArgs([]string{
		"create", "project", ref.Name,
		"--parent", ref.Parent,
		"--symbol", ref.Symbol,
		"--color", fmt.Sprintf("%d", ref.Color),
		"--fg-color", fmt.Sprintf("%d", ref.FgColor),
		"--tags", strings.Join(ref.RequiredTags, ","),
	})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}
	assert.True(t, track.ProjectExists("child"), "Project should exist")

	project, err := track.LoadProject("child")
	if err != nil {
		t.Fatal("error loading project")
	}
	assert.Equal(t, ref, project, "project properties differ from expected due to flags")
}

func TestCreateRecord(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"create", "project", "test"})

	if err := cmd.Execute(); err != nil {
		t.Fatal("error executing command")
	}
	assert.True(t, track.ProjectExists("test"), "Project should exist")

	cmd.SetArgs([]string{"create", "record", "test", "today", "10:00-1h", "note"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("error executing command: %s", err)
	}

	records, err := track.LoadAllRecords()
	if err != nil {
		t.Fatal("error loading records")
	}

	assert.True(t, len(records) == 1, "Record should exist")
	assert.Equal(t, records[0].Note, "note", "Note should be 'note'")

	cmd.SetArgs([]string{"create", "record", "test", "today", "10:00-1h", "note"})
	assert.NotNil(t, cmd.Execute(), "should fail with record time overlap error")

	cmd.SetArgs([]string{"create", "record", "test", "2015-10-21", "16:29-19:28", "back to the future"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("error executing command: %s", err)
	}

	records, err = track.LoadAllRecords()
	if err != nil {
		t.Fatal("error loading records")
	}

	assert.True(t, len(records) == 2, "Records should exist")
	assert.Equal(t, records[0].Note, "back to the future", "Note should be 'back to the future'")
	assert.Equal(t, records[1].Note, "note", "Note should be 'note'")
}
