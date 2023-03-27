package apc

import "errors"

// Returns a parser that parses the exact CT value.
// Returns the result as as CT.
func ExactSlice[CT any](value []CT) Parser[CT, []CT] {
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

// Returns a parser that parses the exact CT value.
// Returns the result as as CT.
func ExactOne[CT any](value CT) Parser[CT, CT] {
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

// Returns a parser that parses the exact string value.
// Returns the result as as string.
func ExactStr(value string) Parser[rune, string] {
	if len(value) <= 0 {
		panic("value for Exact must have a length > 0")
	}

	return func(ctx Context[rune]) (string, error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return "", err
		}

		val, err := ctx.Peek(0, len(value))
		if err != nil && !errors.Is(err, ErrEOF) {
			return "", err
		}
		valStr := string(val)
		if valStr == value {
			_, err := ctx.Consume(len(valStr))
			if err != nil && !errors.Is(err, ErrEOF) {
				return "", err
			}
			return valStr, nil
		}
		return "", ParseErrExpectedButGot(ctx, value, valStr, nil)
	}
}
