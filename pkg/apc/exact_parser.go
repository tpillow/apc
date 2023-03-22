package apc

import (
	"errors"
	"fmt"
)

func Exact(value string) Parser[string] {
	return func(ctx Context) (string, error) {
		debugRunning(fmt.Sprintf("exact %v", value))

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
