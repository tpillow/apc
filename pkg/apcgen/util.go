package apcgen

import (
	"fmt"
	"reflect"
)

// Enable various debug prints
var (
	DebugPrintBuiltNodes = false
)

// Log helper
func maybeLog(doLog bool, format string, formatArgs ...interface{}) {
	if doLog {
		fmt.Printf(fmt.Sprintf("[DEBUG] %v\n", format), formatArgs...)
	}
}

var maybeIntRange = intRange{min: -2, max: -2}

type intRange struct {
	min int
	max int
}

func (ir intRange) String() string {
	return fmt.Sprintf("<IntRange Min=%v Max=%v>", ir.min, ir.max)
}

type keyValuePair[KT, VT any] struct {
	key   KT
	value VT
}

func reflectTypeOf[T any]() reflect.Type {
	var tmp T
	return reflect.TypeOf(tmp)
}
