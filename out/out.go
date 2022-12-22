package out

import (
	"fmt"
)

// Err prints an error
func Err(format string, a ...interface{}) {
	print(format, a...)
}

// Warn prints a warning message
func Warn(format string, a ...interface{}) {
	print(format, a...)
}

// Success prints a success message
func Success(format string, a ...interface{}) {
	print(format, a...)
}

func print(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
