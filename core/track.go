package core

import (
	"path/filepath"

	"github.com/mlange-42/track/fs"
)

// Track is a top-level track instalce
type Track struct {
	Config Config
}

// NewTrack creates a new Track object
func NewTrack() (Track, error) {
	track := Track{}
	track.createRootDir()

	conf, err := LoadConfig()
	if err != nil {
		return track, err
	}

	track.Config = conf
	track.createWorkspaceDirs(track.Config.Workspace)

	return track, nil
}

func (t *Track) createRootDir() {
	err := fs.CreateDir(fs.RootDir())
	if err != nil {
		panic(err)
	}
}

func (t *Track) createWorkspaceDirs(workspace string) {
	err := fs.CreateDir(t.workspaceProjectsDir(workspace))
	if err != nil {
		panic(err)
	}
	err = fs.CreateDir(t.workspaceRecordsDir(workspace))
	if err != nil {
		panic(err)
	}
}

// workspaceProjectsDir returns the projects storage directory for the given workspace
func (t *Track) workspaceProjectsDir(ws string) string {
	return filepath.Join(fs.RootDir(), ws, fs.ProjectsDirName())
}

// workspaceRecordsDir returns the records storage directory for the given workspace
func (t *Track) workspaceRecordsDir(ws string) string {
	return filepath.Join(fs.RootDir(), ws, fs.RecordsDirName())
}
