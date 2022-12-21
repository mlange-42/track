package fs

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	rootDirName     = ".track"
	projectsDirName = "projects"
	recordsDirName  = "records"
)

var pathSanitizer = strings.NewReplacer("/", "-", "\\", "-")

// RootDir returns the root storage directory
func RootDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, rootDirName)
}

// ProjectsDir returns the projects storage directory
func ProjectsDir() string {
	return filepath.Join(RootDir(), projectsDirName)
}

// ProjectFile returns file name for a project
func ProjectFile(name string) string {
	return filepath.Join(ProjectsDir(), pathSanitizer.Replace(name)+".json")
}

// RecordsDir returns the records storage directory
func RecordsDir() string {
	return filepath.Join(RootDir(), recordsDirName)
}
