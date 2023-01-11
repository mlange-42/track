package core

import (
	"fmt"
	"os"

	"github.com/mlange-42/track/util"
)

// CreateWorkspace creates a new workspace
func (t *Track) CreateWorkspace(name string) error {
	if util.DirExists(t.WorkspaceDir(name)) {
		return fmt.Errorf("workspace '%s' already exists", name)
	}
	t.createWorkspaceDirs(name)
	return nil
}

// WorkspaceExists returns whether a workspace exists
func (t *Track) WorkspaceExists(name string) bool {
	return util.DirExists(t.WorkspaceDir(name))
}

// SwitchWorkspace switches to another workspace
func (t *Track) SwitchWorkspace(name string) error {
	if !util.DirExists(t.WorkspaceDir(name)) {
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
	err = t.Config.Save(t.ConfigPath())
	if err != nil {
		return err
	}

	return nil
}

// Workspace returns the current workspace
func (t *Track) Workspace() string {
	return t.Config.Workspace
}

// AllWorkspaces returns a slice of all workspaces
func (t *Track) AllWorkspaces() ([]string, error) {
	path := t.RootDir
	dirs, err := os.ReadDir(path)
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
