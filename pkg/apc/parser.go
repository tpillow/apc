// Package apc provides a minimalist parser combinator library.
package apc

// Parser[T] represents a parser that takes a Context and returns a result of type T or an error.
//
// Should return a nil error if the result was parsed and consumed.
// Should return a ParseError error if parsing failed, and no input was consumed.
// Should return a ParseErrorConsumed error if parsing failed, but some input was consumed.
// Any other error type may be returned, and is treated like ParseErrorConsumed.
type Parser[T any] func(ctx Context) (T, error)

// ParseConfig contains settings that can be passed to the Parse function.
type ParseConfig struct {
	// If true, parsing will fail if there is remaining input in the Context after parsing.
	MustParseToEOF bool
}

// Executes the provided parser using the given context, first applying the parseConfig.
func Parse[T any](ctx Context, parser Parser[T], parseConfig ParseConfig) (T, error) {
	node, err := parser(ctx)
	if err != nil {
		return zeroVal[T](), err
	}

	if parseConfig.MustParseToEOF {
		ctx.RunSkipParsers()
		if _, err := ctx.Peek(0, 1); err == nil {
			return zeroVal[T](), ParseErrExpectedButGotNext(ctx, "EOF", nil)
		}
	}

	return node, nil
}
