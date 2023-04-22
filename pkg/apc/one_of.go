package apc

// Returns a parser that attempts to parse, in order, the provided parsers.
// Returns the result of the first successful parser.
func OneOf[CT, T any](name string, parsers ...Parser[CT, T]) Parser[CT, T] {
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

func LookOneOf[CT, T any](name string, parsers ...Parser[CT, T]) Parser[CT, T] {
	for i := 0; i < len(parsers); i++ {
		parsers[i] = Look(parsers[i])
	}
	return OneOf(name, parsers...)
}
