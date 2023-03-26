package apc

import "unicode/utf8"

// Obtain the zero value of type T.
func zeroVal[T any]() T {
	var val T
	return val
}

// RuneContextPeekingRuneReader implements io.RuneReader using the given
// Context by peeking one rune at a time, starting at offset 0.
type RuneContextPeekingRuneReader struct {
	Context Context[rune]
	offset  int
}

// Peeks the next rune in the Context, and advances the reader offset.
func (r *RuneContextPeekingRuneReader) ReadRune() (rune, int, error) {
	val, err := r.Context.Peek(r.offset, 1)
	if err != nil {
		return 0, 0, err
	}
	r.offset += 1
	rn, size := utf8.DecodeRuneInString(string(val))
	return rn, size, nil
}
