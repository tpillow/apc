package apc

import "errors"

// Returns a parser that succeeds if peeking elements from the Context
// equals value, returning value as the result.
func ExactSlice[CT any](value []CT) Parser[CT, []CT] {
	if len(value) <= 0 {
		panic("value for ExactSlice must have a length > 0")
	}

	return func(ctx Context[CT]) ([]CT, error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return nil, err
		}

		val, err := ctx.Peek(0, len(value))
		if err != nil && !errors.Is(err, ErrEOF) {
			return nil, err
		}
		if len(val) != len(value) {
			return nil, ParseErrExpectedButGot(ctx, value, val, nil)
		}
		for i, r := range value {
			if any(val[i]) != any(r) {
				return nil, ParseErrExpectedButGot(ctx, value, val, nil)
			}
		}
		_, err = ctx.Consume(len(val))
		if err != nil && !errors.Is(err, ErrEOF) {
			return nil, err
		}
		return val, nil
	}
}

// Returns a parser that succeeds if peeking 1 element from the Context
// equals value, returning value as the result.
func Exact[CT any](value CT) Parser[CT, CT] {
	return func(ctx Context[CT]) (CT, error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return zeroVal[CT](), err
		}

		val, err := ctx.Peek(0, 1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return zeroVal[CT](), err
		}
		if any(val[0]) != any(value) {
			return zeroVal[CT](), ParseErrExpectedButGot(ctx, value, val[0], nil)
		}
		_, err = ctx.Consume(1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return zeroVal[CT](), err
		}
		return val[0], nil
	}
}

// Equivalent to ExactSlice but for a Context[rune]. Implicitly converts value to []rune.
func ExactStr(value string) Parser[rune, string] {
	return Map(ExactSlice([]rune(value)), func(node []rune, _ Origin) string {
		return string(node)
	})
}
