package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitize(t *testing.T) {
	tt := []struct {
		title    string
		input    string
		expected string
	}{
		{
			title:    "no illegal characters",
			input:    "test",
			expected: "test",
		},
		{
			title:    "allowed special characters",
			input:    "test$§=()'\".,;-_",
			expected: "test$§=()'\".,;-_",
		},
		{
			title:    "slash and backslash",
			input:    "te/s\\t",
			expected: "te-s-t",
		},
		{
			title:    "braces",
			input:    "t{e}<s>t",
			expected: "t-e--s-t",
		},
		{
			title:    "other",
			input:    "t:e#s~t%t&e*s:t?t+e|st",
			expected: "t-e-s-t-t-e-s-t-t-e-st",
		},
	}

	for _, test := range tt {
		output := Sanitize(test.input)
		assert.Equal(t, test.expected, output, "Wrong sanitized string in %s", test.title)
	}
}

func TestExistsEtc(t *testing.T) {
	dir, err := ioutil.TempDir("", "track-test")
	assert.Nil(t, err, "Error creating temporary directory")
	defer os.Remove(dir)

	assert.True(t, DirExists(dir), "Directory should exist")

	empty, err := DirIsEmpty(dir)
	assert.Nil(t, err, "Error checking empty directory")
	assert.True(t, empty, "Directory should exist")

	assert.False(t, DirExists(filepath.Join(dir, "test")), "Directory should not exist")
	assert.False(t, FileExists(filepath.Join(dir, "test")), "File should not exist")

	file, err := os.Create(filepath.Join(dir, "test.file"))
	assert.Nil(t, err, "Error creating file")
	file.Close()
	assert.True(t, FileExists(filepath.Join(dir, "test.file")), "File should exist")

	empty, err = DirIsEmpty(dir)
	assert.Nil(t, err, "Error checking empty directory")
	assert.False(t, empty, "Directory should exist")

	err = CreateDir(filepath.Join(dir, "foo"))
	assert.Nil(t, err, "Error creating directory")
	assert.True(t, DirExists(filepath.Join(dir, "foo")), "Directory should exist")

	err = CreateDir(filepath.Join(dir, "goo"))
	err = CreateDir(filepath.Join(dir, "hoo"))
	err = CreateDir(filepath.Join(dir, "ioo"))

	path, name, err := FindLatests(dir, true)
	assert.Nil(t, err, "Error finding latest directory")
	assert.Equal(t, filepath.Join(dir, "ioo"), path, "Wrong latest directory path")
	assert.Equal(t, "ioo", name, "Wrong latest directory name")

	file, err = os.Create(filepath.Join(dir, "test2.file"))
	assert.Nil(t, err, "Error creating file")

	path, name, err = FindLatests(dir, false)
	assert.Nil(t, err, "Error finding latest file")
	assert.Equal(t, filepath.Join(dir, "test2.file"), path, "Wrong latest file path")
	assert.Equal(t, "test2.file", name, "Wrong latest file name")
}
