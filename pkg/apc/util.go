package apc

import "unicode/utf8"

func zeroVal[T any]() T {
	var val T
	return val
}

type ContextPeekingRuneReader struct {
	Context Context
	offset  int
}

func (r *ContextPeekingRuneReader) ReadRune() (rune, int, error) {
	val, err := r.Context.Peek(r.offset, 1)
	if err != nil {
		return 0, 0, err
	}
	r.offset += 1
	rn, size := utf8.DecodeRuneInString(val)
	return rn, size, nil
}
