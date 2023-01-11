package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwitchWorkspace(t *testing.T) {
	track, err := setupTestCommand()
	if err != nil {
		t.Fatal("error setting up test")
	}
	defer os.Remove(track.RootDir)

	err = track.CreateWorkspace("move-to")
	if err != nil {
		t.Fatal("error creating workspace")
	}

	cmd := RootCommand(track, "")
	cmd.SetArgs([]string{"workspace", "move-to"})

	err = cmd.Execute()
	if err != nil {
		t.Fatal("error executing command")
	}
	assert.Equal(t, "move-to", track.Workspace(), "Should be in new workspace")
}
