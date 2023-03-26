package apc

import "errors"

type Equatable[T any] interface {
	Equal(other T) bool
}

// Returns a parser that parses the exact CT value.
// Returns the result as as CT.
func Exact[CT Equatable[CT]](value []CT) Parser[CT, []CT] {
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
			if val[i].Equal(r) {
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
