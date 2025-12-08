package apc

import (
	"fmt"
	"regexp"
)

func ExactStr(expectStr string) *Parser[rune, string] {
	if len(expectStr) <= 0 {
		panic("MatchStr requires a string with length > 0")
	}

	desc := fmt.Sprintf("exactly '%s'", expectStr)

	return NewParser(desc, func(ctx Context[rune]) (string, error) {
		runeCtx, ok := ctx.(*sliceContext[rune])
		if !ok {
			panic("MatchStr requires a sliceContext")
		}

		startLoc := runeCtx.location
		remainingLen := len(runeCtx.source) - startLoc.Index
		if remainingLen < len(expectStr) {
			return "", ParseError{
				Up:            nil,
				Expected:      desc,
				Unexpected:    EofString,
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}

		for i := 0; i < len(expectStr); i++ {
			token := ctx.Pop()
			if rune(expectStr[i]) != token {
				endLoc := ctx.CurLocation()
				ctx.SetLocation(startLoc)
				return "", ParseError{
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

func RegexGroup(pattern string, groupIndex int) *Parser[rune, string] {
	regex := regexp.MustCompile("^" + pattern)
	desc := fmt.Sprintf("regex '%s'", pattern)

	return NewParser(desc, func(ctx Context[rune]) (string, error) {
		runeCtx, ok := ctx.(*sliceContext[rune])
		if !ok {
			panic("MatchStr requires a sliceContext")
		}

		startLoc := runeCtx.location
		remaining := string(runeCtx.source[startLoc.Index:])
		matches := regex.FindStringSubmatch(remaining)

		if matches == nil {
			return "", ParseError{
				Up:            nil,
				Expected:      desc,
				Unexpected:    "unmatched regex",
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}
		fullMatchGroup := matches[0]
		desiredMatchGroup := matches[groupIndex]
		if len(fullMatchGroup) <= 0 {
			panic("a Regex Parser cannot have a 0 length match")
		}

		for i := 0; i < len(fullMatchGroup); i++ {
			ctx.Pop()
		}

		return desiredMatchGroup, nil
	})
}

func Regex(pattern string) *Parser[rune, string] {
	return RegexGroup(pattern, 0)
}
