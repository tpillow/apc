package apcgen

import "fmt"

var maybeIntRange = IntRange{Min: -2, Max: -2}

type IntRange struct {
	Min int
	Max int
}

func (ir IntRange) String() string {
	return fmt.Sprintf("<IntRange Min=%v Max=%v>", ir.Min, ir.Max)
}

type keyValuePair[KT, VT any] struct {
	key   KT
	value VT
}
