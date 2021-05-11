package twig

import "fmt"

var debug = false

// Switch changes the debug stats
func Switch() {
	debug = !debug
	if debug {
		Printf("debug ON")
	}
}

// Printf writes to the standard output if the debug is enabled
func Printf(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format+"\n", args...)
	}
}
