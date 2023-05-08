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

func anyConvertRunesToString(rawVal any) string {
	switch val := rawVal.(type) {
	case rune:
		return string(val)
	case []rune:
		return string(val)
	default:
		return fmt.Sprintf("%v", val)
	}
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
// Converts any "" or "[]" value to "EOF".
func interfaceToErrString(val interface{}) string {
	// TODO: a better way???
	ret := fmt.Sprintf("%v", val)
	if cval, ok := val.([]int32); ok {
		ret = string(cval)
	}
	if ret == "" || ret == "[]" {
		return "EOF"
	}
	return ret
}

// Represents any value along with if the value is nil. Useful when working with non-pointer types.
type MaybeValue[T any] struct {
	isNil bool
	value T
}

// Returns true if the value represented is nil.
func (val MaybeValue[T]) IsNil() bool {
	return val.isNil
}

// Returns the value represented. Panics if IsNil().
func (val MaybeValue[T]) Value() T {
	if val.isNil {
		panic("cannot call Value() on a nil MaybeValue")
	}
	return val.value
}

// Returns a MaybeValue that returns true for IsNil().
func NewNilMaybeValue[T any]() MaybeValue[T] {
	return MaybeValue[T]{isNil: true}
}

// Returns a MaybeValue that returns the provide value for Value().
func NewMaybeValue[T any](value T) MaybeValue[T] {
	return MaybeValue[T]{isNil: false, value: value}
}
