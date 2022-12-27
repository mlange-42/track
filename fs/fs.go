package fs

import (
	"errors"
	gofs "io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	rootDirName     = ".track"
	projectsDirName = "projects"
	recordsDirName  = "records"
	configFile      = "config.yml"
)

var (
	// ErrNoFiles is an error for no files found at all
	ErrNoFiles = errors.New("no files")
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

// ProjectsDirName returns the directory name for projects
func ProjectsDirName() string {
	return projectsDirName
}

// RecordsDirName returns the directory name for records
func RecordsDirName() string {
	return recordsDirName
}

// ConfigPath returns the default config path
func ConfigPath() string {
	return filepath.Join(RootDir(), configFile)
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

// FindLatests finds the "latest" file or directory in a file, by name
func FindLatests(path string, isDir bool) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	var dir gofs.FileInfo = nil
	for i := len(files) - 1; i >= 0; i-- {
		dir = files[i]
		if dir.IsDir() == isDir {
			break
		}
	}

	if dir == nil {
		return "", ErrNoFiles
	}

	return filepath.Join(path, dir.Name()), nil
}
