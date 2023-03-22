package apc

import "errors"

func OneOf[T any](name string, parsers ...Parser[T]) Parser[T] {
	if len(parsers) < 2 {
		panic("must provide at least 2 parsers to OneOf")
	}

	return func(ctx Context) (T, error) {
		for _, parser := range parsers {
			node, err := parser(ctx)
			if err == nil {
				return node, nil
			}
			if !errors.Is(err, ErrParseErr) {
				return zeroVal[T](), err
			}
		}
		return zeroVal[T](), ParseErrExpectedButGotNext(ctx, name, nil)
	}
}