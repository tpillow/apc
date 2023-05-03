package apc

import "fmt"

// Origin represents a line and column number location in some source.
type Origin struct {
	// The name of the origin (usually filename).
	Name string
	// The line number location.
	LineNum int
	// The column number location.
	ColNum int
}

// Returns a string representation of an Origin.
func (origin Origin) String() string {
	return fmt.Sprintf("%v:%v:%v", origin.Name, origin.LineNum, origin.ColNum)
}

// Holds a start and end Origin.
type OriginRange struct {
	// The starting origin.
	Start Origin
	// The end origin.
	End Origin
}
