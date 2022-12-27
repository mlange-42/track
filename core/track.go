package core

import (
	"fmt"

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
	track.createWorkspaceDirs()

	return track, nil
}

// Workspace returns the current workspace
func (t *Track) Workspace() string {
	return t.Config.Workspace
}

// WorkspaceLabel returns the current workspace label
func (t *Track) WorkspaceLabel() string {
	return fmt.Sprintf(RootPattern, t.Config.Workspace)
}

func (t *Track) createRootDir() {
	err := fs.CreateDir(fs.RootDir())
	if err != nil {
		panic(err)
	}
}

func (t *Track) createWorkspaceDirs() {
	err := fs.CreateDir(t.ProjectsDir())
	if err != nil {
		panic(err)
	}
	err = fs.CreateDir(t.RecordsDir())
	if err != nil {
		panic(err)
	}
}
