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
)

var (
	// ErrNoFiles is an error for files found at all
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
