package out

import (
	"fmt"
	"os"
)

// WriteToStdout writes a message to standard output
func WriteToStdout(message string) {
	fmt.Fprintln(os.Stdout, message)
}

// WriteToStderr writes a message to standard error
func WriteToStderr(message string) {
	fmt.Fprintln(os.Stderr, message)
}

// ExitSuccess exits the program with status code 0 (success)
func ExitSuccess() {
	os.Exit(0)
}

// ExitError exits the program with status code 1 (error) and prints an error message
func ExitError(message string) {
	WriteToStderr(message)
	os.Exit(1)
}
