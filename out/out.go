package out

import (
	"fmt"
	"io"
	"os"

	"github.com/gookit/color"
)

var (
	successColor = color.BgGreen
	warningColor = color.BgYellow
	errorColor   = color.BgRed
	promptColor  = color.BgBlue
)

// StdOut is the writer for standard output. Defaults to os.Stdout
var StdOut io.Writer = os.Stdout

// StdErr is the writer for error output. Defaults to os.Stderr
var StdErr io.Writer = os.Stderr

// StdIn is the reader for standard input. Defaults to os.Stdin
var StdIn io.Reader = os.Stdin

// Print prints a neutral message
func Print(format string, a ...interface{}) {
	printOut(format, a...)
}

// Err prints an error
func Err(format string, a ...interface{}) {
	fmt.Fprintf(StdOut, "%s ", errorColor.Sprint(" ERROR   "))
	printErr(format, a...)
}

// Warn prints a warning message
func Warn(format string, a ...interface{}) {
	fmt.Fprintf(StdOut, "%s ", warningColor.Sprint(" WARNING "))
	printErr(format, a...)
}

// Success prints a success message
func Success(format string, a ...interface{}) {
	fmt.Fprintf(StdOut, "%s ", successColor.Sprint(" SUCCESS "))
	printErr(format, a...)
}

// Scan prints a prompt message and scans for user input
func Scan(format string, a ...interface{}) (string, error) {
	fmt.Fprintf(StdOut, "%s ", promptColor.Sprint(" PROMPT  "))
	printOut(format, a...)

	var answer string
	_, err := fmt.Fscanln(StdIn, &answer)
	return answer, err
}

func printOut(format string, a ...interface{}) {
	fmt.Fprintf(StdOut, format, a...)
}

func printErr(format string, a ...interface{}) {
	fmt.Fprintf(StdErr, format, a...)
}
