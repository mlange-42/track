package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mlange-42/track/fs"
)

// Project holds and manipulates data for a project
type Project struct {
	Name   string
	Parent string
}

// ProjectFile returns file name for a project
func (t *Track) ProjectFile(name string) string {
	return filepath.Join(fs.ProjectsDir(), fs.Sanitize(name)+".json")
}

// ProjectExists checks if a project exists
func (t *Track) ProjectExists(name string) bool {
	return fs.FileExists(t.ProjectFile(name))
}

// SaveProject saves a project to disk
func (t *Track) SaveProject(project Project) error {
	path := t.ProjectFile(project.Name)

	if fs.FileExists(path) {
		return fmt.Errorf("Project '%s' already exists", project.Name)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(&project, "", "\t")
	if err != nil {
		return err
	}

	_, err = file.Write(bytes)

	return err
}
