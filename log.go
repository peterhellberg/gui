package gui

import "fmt"

// Log to standard output.
func Log(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}
