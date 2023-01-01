package out

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

var (
	successColor = color.BgGreen
	warningColor = color.BgYellow
	errorColor   = color.BgRed
	promptColor  = color.BgBlue
)

// Print prints a neutral message
func Print(format string, a ...interface{}) {
	printOut(format, a...)
}

// Err prints an error
func Err(format string, a ...interface{}) {
	errorColor.Print(" ERROR   ")
	fmt.Print(" ")
	printErr(format, a...)
}

// Warn prints a warning message
func Warn(format string, a ...interface{}) {
	warningColor.Print(" WARNING ")
	fmt.Print(" ")
	printErr(format, a...)
}

// Success prints a success message
func Success(format string, a ...interface{}) {
	successColor.Print(" SUCCESS ")
	fmt.Print(" ")
	printErr(format, a...)
}

// Scan prints a prompt message and scans foruser input
func Scan(format string, a ...interface{}) (string, error) {
	promptColor.Print(" PROMPT  ")
	fmt.Print(" ")
	printOut(format, a...)

	var answer string
	_, err := fmt.Scanln(&answer)
	return answer, err
}

func printOut(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func printErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
