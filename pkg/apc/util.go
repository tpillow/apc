package apc

import (
	"fmt"
	"unicode/utf8"
)

// Enable various debug prints
var (
	DebugPrintReaderContext = false
)

// Log helper
func maybeLog(doLog bool, format string, formatArgs ...interface{}) {
	if doLog {
		fmt.Printf(fmt.Sprintf("[DEBUG] %v\n", format), formatArgs...)
	}
}

// Obtain the zero value of type T.
func zeroVal[T any]() T {
	var val T
	return val
}

// RuneContextPeekingRuneReader implements io.RuneReader using the given
// Context[rune] by peeking one rune at a time, starting at offset 0.
type RuneContextPeekingRuneReader struct {
	Context Context[rune]
	offset  int
}

// Peeks the next rune in the Context[rune], and advances the reader offset.
func (r *RuneContextPeekingRuneReader) ReadRune() (rune, int, error) {
	val, err := r.Context.Peek(r.offset, 1)
	if err != nil {
		return 0, 0, err
	}
	r.offset += 1
	rn, size := utf8.DecodeRuneInString(string(val))
	return rn, size, nil
}

// Turns a generic interface{} into a string that is sufficient to report in an error message.
// Converts any []int32 to a string.
// Converts any "" or "[]" value to "<nothing>".
func interfaceToErrString(val interface{}) string {
	// TODO: a better way???
	ret := fmt.Sprintf("%v", val)
	if cval, ok := val.([]int32); ok {
		ret = string(cval)
	}
	if ret == "" || ret == "[]" {
		return "<nothing>"
	}
	return ret
}
