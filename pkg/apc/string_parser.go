package apc

import (
	"fmt"
	"regexp"
)

func ExactStr(expectStr string) Parser {
	if len(expectStr) <= 0 {
		panic("MatchStr requires a string with length > 0")
	}

	desc := fmt.Sprintf("exactly '%s'", expectStr)

	return NewParser(desc, func(ctx Context) (any, error) {
		runeCtx, ok := ctx.(*runeSliceContext)
		if !ok {
			panic("MatchStr requires a sliceContext")
		}

		startLoc := runeCtx.location
		remainingLen := len(runeCtx.source) - startLoc.Index
		if remainingLen < len(expectStr) {
			return nil, ParseError{
				Up:            nil,
				Expected:      desc,
				Unexpected:    EofToken{},
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}

		for i := 0; i < len(expectStr); i++ {
			token := ctx.Pop()
			if expectStr[i] != token {
				endLoc := ctx.CurLocation()
				ctx.SetLocation(startLoc)
				return nil, ParseError{
					Up:            nil,
					Expected:      desc,
					Unexpected:    token,
					StartLocation: startLoc,
					EndLocation:   endLoc,
				}
			}
		}

		return expectStr, nil
	})
}

func RegexGroup(description string, pattern string, groupIndex int) Parser {
	regex := regexp.MustCompile("^" + pattern)

	return NewParser(description, func(ctx Context) (any, error) {
		runeCtx, ok := ctx.(*runeSliceContext)
		if !ok {
			panic("MatchStr requires a sliceContext")
		}

		startLoc := runeCtx.sliceContext.location
		remaining := string(runeCtx.sourceAsRuneSlice[startLoc.Index:])
		matches := regex.FindStringSubmatch(remaining)

		if matches == nil {
			return "", ParseError{
				Up:            nil,
				Expected:      "regex to match",
				Unexpected:    "unmatched regex",
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}
		matchGroup := matches[groupIndex]
		if len(matchGroup) <= 0 {
			panic("a Regex Parser cannot have a 0 length match")
		}

		for i := 0; i < len(matchGroup); i++ {
			if IsEofToken(ctx.Pop()) {
				panic("unreachable")
			}
		}

		return matchGroup, nil
	})
}

func Regex(description string, pattern string) Parser {
	return RegexGroup(description, pattern, 0)
}
