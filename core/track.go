package core

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/mlange-42/track/fs"
)

// Track is a top-level track instalce
type Track struct {
}

// Project holds and manipulates data for a project
type Project struct {
	Name   string
	Parent string
}

// Record holds and manipulates data for a record
type Record struct {
	Project Project
	Note    string
	Start   time.Time
	End     time.Time
}

// SaveProject saves a project to disk
func (t *Track) SaveProject(project Project) error {
	path := fs.ProjectFile(project.Name)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("Project '%s' already exists", project.Name)
	}

	bytes, err := json.MarshalIndent(&project, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}

// CreateDirs creates the storage directories
func (t *Track) CreateDirs() {
	err := createDir(fs.ProjectsDir())
	if err != nil {
		panic(err)
	}
	err = createDir(fs.RecordsDir())
	if err != nil {
		panic(err)
	}
}

func createDir(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return nil
	}
	return err
}
