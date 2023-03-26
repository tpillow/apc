package apc

import "errors"

// Returns a parser that parses the exact string value.
// Returns the result as as string.
func ExactStr(value string) Parser[string, string] {
	if len(value) <= 0 {
		panic("value for Exact must have a length > 0")
	}

	return func(ctx Context[string]) (string, error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return "", err
		}

		val, err := ctx.Peek(0, len(value))
		if err != nil && !errors.Is(err, ErrEOF) {
			return "", err
		}
		if val == value {
			_, err := ctx.Consume(len(val))
			if err != nil && !errors.Is(err, ErrEOF) {
				return "", err
			}
			return val, nil
		}
		return "", ParseErrExpectedButGot(ctx, value, val, nil)
	}
}
