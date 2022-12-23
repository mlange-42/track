package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mlange-42/track/fs"
)

// Project holds and manipulates data for a project
type Project struct {
	Name   string
	Parent string
}

// ProjectPath returns the full path for a project
func (t *Track) ProjectPath(name string) string {
	return filepath.Join(fs.ProjectsDir(), fs.Sanitize(name)+".json")
}

// ProjectExists checks if a project exists
func (t *Track) ProjectExists(name string) bool {
	return fs.FileExists(t.ProjectPath(name))
}

// SaveProject saves a project to disk
func (t *Track) SaveProject(project Project, force bool) error {
	path := t.ProjectPath(project.Name)

	if !force && fs.FileExists(path) {
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

// LoadProjectByName loads a project by name
func (t *Track) LoadProjectByName(name string) (Project, error) {
	path := t.ProjectPath(name)
	return t.LoadProject(path)
}

// LoadProject loads a project
func (t *Track) LoadProject(path string) (Project, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Project{}, err
	}

	var project Project

	if err := json.Unmarshal(file, &project); err != nil {
		return Project{}, err
	}

	return project, nil
}

// LoadAllProjects loads all projects
func (t *Track) LoadAllProjects() (map[string]Project, error) {
	path := fs.ProjectsDir()

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	projects := make(map[string]Project)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		project, err := t.LoadProject(filepath.Join(path, file.Name()))
		if err != nil {
			return nil, err
		}
		projects[project.Name] = project
	}

	return projects, nil
}
