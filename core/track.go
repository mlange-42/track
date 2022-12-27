package core

import (
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
	err := fs.CreateDir(t.WorkspaceProjectsDir(workspace))
	if err != nil {
		panic(err)
	}
	err = fs.CreateDir(t.WorkspaceRecordsDir(workspace))
	if err != nil {
		panic(err)
	}
}
