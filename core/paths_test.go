package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPaths(t *testing.T) {
	home, err := os.UserHomeDir()
	assert.Nil(t, err, "Error getting home directory")

	track, err := NewTrack(nil)

	assert.Nil(t, err, "Error initializing Track")
	assert.Equal(t, filepath.Join(home, rootDirName), track.RootDir, "Wrong root directory")
	assert.Equal(t, filepath.Join(home, rootDirName, configFile), track.ConfigPath(), "Wrong config path")
	assert.Equal(t, filepath.Join(home, rootDirName, defaultWorkspace), track.WorkspaceDir(defaultWorkspace), "Wrong workspace directory")
	assert.Equal(t, filepath.Join(home, rootDirName, defaultWorkspace, projectsDirName), track.ProjectsDir(), "Wrong projects directory")
	assert.Equal(t, filepath.Join(home, rootDirName, defaultWorkspace, recordsDirName), track.RecordsDir(), "Wrong records directory")

	projects := track.ProjectsDir()
	records := track.RecordsDir()

	assert.Equal(t, filepath.Join(projects, "test.yml"), track.ProjectPath("test"), "Wrong project file")
	assert.Equal(t, filepath.Join(records, "2001", "02", "03"), track.RecordDir(time.Date(2001, 2, 3, 4, 5, 0, 0, time.Local)), "Wrong record directory")
	assert.Equal(t, filepath.Join(records, "2001", "02", "03", "04-05.trk"), track.RecordPath(time.Date(2001, 2, 3, 4, 5, 0, 0, time.Local)), "Wrong record file")
}
