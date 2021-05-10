package log

import "fmt"

var debug = false

func Switch() {
	debug = !debug
	if debug {
		Log("debug ON")
	}
}

func Log(format string, args ...interface{}) {
	if debug {
		fmt.Printf(format+"\n", args...)
	}
}
