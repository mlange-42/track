package core

import (
	"github.com/mlange-42/track/fs"
)

// Track is a top-level track instalce
type Track struct {
}

// CreateDirs creates the storage directories
func (t *Track) CreateDirs() {
	err := fs.CreateDir(fs.ProjectsDir())
	if err != nil {
		panic(err)
	}
	err = fs.CreateDir(fs.RecordsDir())
	if err != nil {
		panic(err)
	}
}
