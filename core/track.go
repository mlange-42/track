package core

import (
	"os"
	"path/filepath"

	"github.com/mlange-42/track/fs"
)

const (
	rootDirName     = ".track"
	projectsDirName = "projects"
	recordsDirName  = "records"
	configFile      = "config.yml"
	trackPathEnvVar = "TRACK_PATH"
)

// Track is a top-level track instalce
type Track struct {
	RootDir string
	Config  Config
}

// NewTrack creates a new Track object
func NewTrack(root *string) (Track, error) {
	track := Track{
		RootDir: getRootDir(root),
	}
	track.createRootDir()

	conf, err := LoadConfig(track.ConfigPath())
	if err != nil {
		return track, err
	}

	track.Config = conf
	track.createWorkspaceDirs(track.Config.Workspace)

	return track, nil
}

func getRootDir(root *string) string {
	if root != nil {
		return *root
	}
	if path, ok := os.LookupEnv(trackPathEnvVar); ok {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, rootDirName)
}

func (t *Track) createRootDir() {
	err := fs.CreateDir(t.RootDir)
	if err != nil {
		panic(err)
	}
}

func (t *Track) createWorkspaceDirs(workspace string) {
	err := fs.CreateDir(t.workspaceProjectsDir(workspace))
	if err != nil {
		panic(err)
	}
	err = fs.CreateDir(t.workspaceRecordsDir(workspace))
	if err != nil {
		panic(err)
	}
}

// workspaceProjectsDir returns the projects storage directory for the given workspace
func (t *Track) workspaceProjectsDir(ws string) string {
	return filepath.Join(t.RootDir, ws, t.ProjectsDirName())
}

// workspaceRecordsDir returns the records storage directory for the given workspace
func (t *Track) workspaceRecordsDir(ws string) string {
	return filepath.Join(t.RootDir, ws, t.RecordsDirName())
}
