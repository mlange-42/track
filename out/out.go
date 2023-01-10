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
	fmt.Fprintf(os.Stdout, "%s ", errorColor.Sprint(" ERROR   "))
	printErr(format, a...)
}

// Warn prints a warning message
func Warn(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, "%s ", warningColor.Sprint(" WARNING "))
	printErr(format, a...)
}

// Success prints a success message
func Success(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, "%s ", successColor.Sprint(" SUCCESS "))
	printErr(format, a...)
}

// Scan prints a prompt message and scans for user input
func Scan(format string, a ...interface{}) (string, error) {
	fmt.Fprintf(os.Stdout, "%s ", promptColor.Sprint(" PROMPT  "))
	printOut(format, a...)

	var answer string
	_, err := fmt.Scanln(&answer)
	return answer, err
}

func printOut(format string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, format, a...)
}

func printErr(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
