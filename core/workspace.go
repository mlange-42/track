package core

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/mlange-42/track/fs"
)

// CreateWorkspace creates a new workspace
func (t *Track) CreateWorkspace(name string) error {
	if fs.DirExists(t.WorkspaceDir(name)) {
		return fmt.Errorf("workspace '%s' already exists", name)
	}
	t.createWorkspaceDirs(name)
	return nil
}

// SwitchWorkspace switches to another workspace
func (t *Track) SwitchWorkspace(name string) error {
	if !fs.DirExists(t.WorkspaceDir(name)) {
		return fmt.Errorf("workspace '%s' does not exist", name)
	}
	open, err := t.OpenRecord()
	if err != nil {
		return err
	}
	if open != nil {
		return fmt.Errorf("running record in workspace '%s'", t.Workspace())
	}
	t.createWorkspaceDirs(name)

	t.Config.Workspace = name
	SaveConfig(t.Config)

	return nil
}

// WorkspaceDir returns the directory of a workspace
func (t *Track) WorkspaceDir(ws string) string {
	return filepath.Join(fs.RootDir(), ws)
}

// Workspace returns the current workspace
func (t *Track) Workspace() string {
	return t.Config.Workspace
}

// AllWorkspaces returns a slice of all workspaces
func (t *Track) AllWorkspaces() ([]string, error) {
	path := fs.RootDir()
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, f := range dirs {
		if !f.IsDir() {
			continue
		}
		result = append(result, f.Name())
	}
	return result, nil
}

// WorkspaceLabel returns the current workspace label
func (t *Track) WorkspaceLabel() string {
	return fmt.Sprintf(RootPattern, t.Config.Workspace)
}
