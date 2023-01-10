package core

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspace(t *testing.T) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(t, err, "Error creating temporary directory")
	defer os.Remove(dir)

	track, err := NewTrack(&dir)
	assert.Nil(t, err, "Error creating Track instance")

	assert.False(t, track.WorkspaceExists("foo"), "Workspace foo should not exist")
	assert.True(t, track.WorkspaceExists("default"), "Workspace default should exist")

	assert.Equal(t, "default", track.Workspace(), "Default workspace should be default")
	assert.Equal(t, "<default>", track.WorkspaceLabel(), "Default workspace label should be <default>")

	err = track.CreateWorkspace("test-ws")
	assert.Nil(t, err, "Error creating workspace")
	assert.True(t, track.WorkspaceExists("test-ws"), "Workspace test-ws should exist")

	err = track.SwitchWorkspace("test-ws")
	assert.Nil(t, err, "Error switching workspace")
	assert.Equal(t, "test-ws", track.Workspace(), "Workspace should be test-ws")

	allWs, err := track.AllWorkspaces()
	assert.Nil(t, err, "Error listing workspace")
	assert.Equal(t, []string{"default", "test-ws"}, allWs, "Workspace should be test-ws")
}
