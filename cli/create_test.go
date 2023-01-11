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

	project, err := track.LoadProjectByName("child")
	if err != nil {
		t.Fatal("error loading project")
	}
	assert.Equal(t, ref, project, "project properties differ from expected due to flags")
}
