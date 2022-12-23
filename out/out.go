package out

import (
	"fmt"

	"github.com/gookit/color"
)

var (
	successColor = color.BgGreen
	warningColor = color.BgYellow
	errorColor   = color.BgRed
)

// Print prints a neutral message
func Print(format string, a ...interface{}) {
	print(format, a...)
}

// Err prints an error
func Err(format string, a ...interface{}) {
	errorColor.Print(" ERROR ")
	fmt.Print(" ")
	print(format, a...)
}

// Warn prints a warning message
func Warn(format string, a ...interface{}) {
	warningColor.Print(" WARNING ")
	fmt.Print(" ")
	print(format, a...)
}

// Success prints a success message
func Success(format string, a ...interface{}) {
	successColor.Print(" SUCCESS ")
	fmt.Print(" ")
	print(format, a...)
}

func print(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
