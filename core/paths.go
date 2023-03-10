package core

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mlange-42/track/util"
)

// WorkspaceDir returns the directory of a workspace
func (t *Track) WorkspaceDir(ws string) string {
	return filepath.Join(t.RootDir, ws)
}

// ConfigPath returns the default config path
func (t *Track) ConfigPath() string {
	return filepath.Join(t.RootDir, configFile)
}

// ProjectsDirName returns the directory name for projects
func (t *Track) ProjectsDirName() string {
	return projectsDirName
}

// ProjectsDir returns the projects storage directory
func (t *Track) ProjectsDir() string {
	return filepath.Join(t.RootDir, t.Workspace(), t.ProjectsDirName())
}

// ProjectPath returns the full path for a project
func (t *Track) ProjectPath(name string) string {
	return filepath.Join(t.ProjectsDir(), util.Sanitize(name)+".yml")
}

// RecordsDirName returns the directory name for records
func (t *Track) RecordsDirName() string {
	return recordsDirName
}

// RecordDir returns the directory path for a record
func (t *Track) RecordDir(tm time.Time) string {
	return filepath.Join(
		t.RecordsDir(),
		fmt.Sprintf("%04d", tm.Year()),
		fmt.Sprintf("%02d", int(tm.Month())),
		fmt.Sprintf("%02d", tm.Day()),
	)
}

// RecordsDir returns the records storage directory
func (t *Track) RecordsDir() string {
	return filepath.Join(t.RootDir, t.Workspace(), t.RecordsDirName())
}

// RecordPath returns the full path for a record
func (t *Track) RecordPath(tm time.Time) string {
	return filepath.Join(
		t.RecordDir(tm),
		fmt.Sprintf("%s.trk", tm.Format(util.FileTimeFormat)),
	)
}
