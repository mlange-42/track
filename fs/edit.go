package fs

import (
	"os"
	"os/exec"
)

// EditFile opens a file in the default editor and waits for the process to finish
func EditFile(path string, editor string) error {
	cmd := exec.Command(editor, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
