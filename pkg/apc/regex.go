package apc

import (
	"errors"
	"fmt"
	"regexp"
)

// Returns a parser that succeeds if peeking elements from the Context[rune]
// matches the provided regex pattern, returning the match as the result.
//
// Note that the regex is always normalized to contain '^' as the starting
// symbol, to always match the left-most character in the input stream.
func Regex(pattern string) Parser[rune, string] {
	if len(pattern) < 1 {
		panic("regex pattern length must be >= 1")
	}
	if pattern[0] != '^' {
		pattern = fmt.Sprintf("^%v", pattern)
	}
	regex := regexp.MustCompile(pattern)

	return func(ctx Context[rune]) (string, error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return "", err
		}

		reader := &RuneContextPeekingRuneReader{Context: ctx}
		loc := regex.FindReaderIndex(reader)
		if loc == nil {
			return "", ParseErrExpectedButGotNext(ctx, ctx.PeekName(), nil)
		}
		if loc[0] != 0 {
			panic("regex should always be normalized to match at start of line")
		}

		matchVal, err := ctx.Consume(loc[1])
		if err != nil && !errors.Is(err, ErrEOF) {
			return "", err
		}
		return string(matchVal), nil
	}
}
