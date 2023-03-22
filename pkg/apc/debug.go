package apc

import "fmt"

var DebugPrint bool = false

func debugRunning(name string) {
	if !DebugPrint {
		return
	}
	debug("running %v\n", name)
}

func debug(format string, args ...interface{}) {
	if !DebugPrint {
		return
	}
	fmt.Printf("[DEBUG] %v", fmt.Sprintf(format, args...))
}
