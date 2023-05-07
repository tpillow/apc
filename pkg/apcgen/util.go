package apcgen

import "fmt"

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
