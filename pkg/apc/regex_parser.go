package apc

import (
	"errors"
	"fmt"
	"regexp"
)

const regexPeekBufferSize = 1024 // TODO: better way (can make RuneReader wrapped around Peek)

func Regex(name string, pattern string) Parser[string] {
	if pattern[0] != '^' {
		pattern = fmt.Sprintf("^%v", pattern)
	}
	regex := regexp.MustCompile(pattern)

	return func(ctx Context) (string, error) {
		debugRunning(name)

		val, err := ctx.Peek(0, regexPeekBufferSize)
		if err != nil && !errors.Is(err, ErrEOF) {
			return "", err
		}
		loc := regex.FindStringIndex(val)
		if loc == nil {
			return "", ParseErrExpectedButGotNext(ctx, name, nil)
		}
		if loc[0] != 0 {
			panic("regex should always be normalized to match at start of line")
		}
		matchVal := val[:loc[1]]
		_, err = ctx.Consume(len(matchVal))
		if err != nil && !errors.Is(err, ErrEOF) {
			return "", err
		}
		return matchVal, nil
	}
}
