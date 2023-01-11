package util

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

var (
	// ErrNoFiles is an error for no files found at all
	ErrNoFiles = errors.New("no files")
)

var illegalFileCharacters = regexp.MustCompile(`[\~\#\%\&\*\{\}\\\:\<\>\?\/\+\|]`)

// Sanitize makes stings filename compatible
func Sanitize(file string) string {
	return illegalFileCharacters.ReplaceAllString(file, "-")
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

// DirIsEmpty checks if a directory is empty
func DirIsEmpty(path string) (bool, error) {
	if !DirExists(path) {
		return false, fmt.Errorf("is not a directory: %s", path)
	}
	content, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(content) == 0, nil
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
func FindLatests(path string, isDir bool) (string, string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return "", "", err
	}

	var dir fs.DirEntry = nil
	for i := len(files) - 1; i >= 0; i-- {
		dir = files[i]
		if dir.IsDir() == isDir {
			break
		}
	}

	if dir == nil {
		return "", "", ErrNoFiles
	}

	return filepath.Join(path, dir.Name()), dir.Name(), nil
}
