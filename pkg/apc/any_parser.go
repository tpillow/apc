package apc

import "errors"

func Any[T any](name string, parsers ...Parser[T]) Parser[T] {
	if len(parsers) < 2 {
		panic("must provide at least 2 parsers to Any")
	}

	return func(ctx Context) (T, error) {
		ctx.RunSkipParsers()

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
