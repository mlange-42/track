package util

import (
	"os"
	"os/exec"
)

// SkipEditingForTests makes editing just fall through, for unit testing
var SkipEditingForTests = false

// EditFile opens a file in the default editor and waits for the process to finish
func EditFile(path string, editor string) error {
	if SkipEditingForTests {
		return nil
	}
	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
