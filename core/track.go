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
	track.createDirs()

	conf, err := LoadConfig()
	if err != nil {
		return track, err
	}

	track.Config = conf
	return track, nil
}

// createDirs creates the storage directories
func (t *Track) createDirs() {
	err := fs.CreateDir(fs.ProjectsDir())
	if err != nil {
		panic(err)
	}
	err = fs.CreateDir(fs.RecordsDir())
	if err != nil {
		panic(err)
	}
}
