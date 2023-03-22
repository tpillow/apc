package apc

import (
	"errors"
	"fmt"
	"regexp"
)

func Regex(name string, pattern string) Parser[string] {
	if len(pattern) < 1 {
		panic("regex pattern length must be >= 1")
	}
	if pattern[0] != '^' {
		pattern = fmt.Sprintf("^%v", pattern)
	}
	regex := regexp.MustCompile(pattern)

	return func(ctx Context) (string, error) {
		ctx.RunSkipParsers()

		reader := &ContextPeekingRuneReader{Context: ctx}
		loc := regex.FindReaderIndex(reader)
		if loc == nil {
			return "", ParseErrExpectedButGotNext(ctx, name, nil)
		}
		if loc[0] != 0 {
			panic("regex should always be normalized to match at start of line")
		}

		matchVal, err := ctx.Consume(loc[1])
		if err != nil && !errors.Is(err, ErrEOF) {
			return "", err
		}
		return matchVal, nil
	}
}
