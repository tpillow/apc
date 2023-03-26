package apc

// Returns a parser that attempts to parse, in order, the provided parsers.
// Returns the result of the first successful parser.
func OneOf[CT any, T any](name string, parsers ...Parser[CT, T]) Parser[CT, T] {
	if len(parsers) < 2 {
		panic("must provide at least 2 parsers to OneOf")
	}

	return func(ctx Context[CT]) (T, error) {
		for _, parser := range parsers {
			node, err := parser(ctx)
			if err == nil {
				return node, nil
			}
			if IsMustReturnParseErr(err) {
				return zeroVal[T](), err
			}
		}
		return zeroVal[T](), ParseErrExpectedButGotNext(ctx, name, nil)
	}
}
