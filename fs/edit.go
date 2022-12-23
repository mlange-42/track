package fs

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// EditFile opens a file in the default editor and waits for the process to finish
func EditFile(path string) error {
	var prog string
	switch os := strings.ToLower(runtime.GOOS); os {
	case "windows":
		prog = "notepad.exe"
	case "linux", "darwin":
		prog = "nano"
	default:
		return fmt.Errorf("unsupported OS: %s", os)
	}
	cmd := exec.Command(prog, path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
