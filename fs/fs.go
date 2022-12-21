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

// RecordsDir returns the records storage directory
func RecordsDir() string {
	return filepath.Join(RootDir(), recordsDirName)
}

// Sanitize makes stings filename compatible
func Sanitize(file string) string {
	return pathSanitizer.Replace(file)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CreateDir creates directories recursively
func CreateDir(path string) error {
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
