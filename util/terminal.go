package util

import (
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// TerminalSize returns the size of the terminal
func TerminalSize() (width int, height int, err error) {
	return terminal.GetSize(int(os.Stdout.Fd()))
}
