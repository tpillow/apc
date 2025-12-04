package apc

import (
	"fmt"
	"regexp"
)

func ExactStr(expectStr string) Parser {
	if len(expectStr) <= 0 {
		panic("MatchStr requires a string with length > 0")
	}

	return NewParser(func(ctx Context) (any, error) {
		runeCtx, ok := ctx.(*runeSliceContext)
		if !ok {
			panic("MatchStr requires a sliceContext")
		}

		startLoc := runeCtx.location
		remainingLen := len(runeCtx.source) - startLoc.Index
		if remainingLen < len(expectStr) {
			return nil, fmt.Errorf("unexpected EOF")
		}

		for i := 0; i < len(expectStr); i++ {
			token := ctx.Pop()
			if expectStr[i] != token {
				ctx.SetLocation(startLoc)
				return nil, fmt.Errorf("not matched")
			}
		}

		return expectStr, nil
	})
}

func RegexGroup(pattern string, groupIndex int) Parser {
	regex := regexp.MustCompile("^" + pattern)

	return NewParser(func(ctx Context) (any, error) {
		runeCtx, ok := ctx.(*runeSliceContext)
		if !ok {
			panic("MatchStr requires a sliceContext")
		}

		startLoc := runeCtx.sliceContext.location
		remaining := string(runeCtx.sourceAsRuneSlice[startLoc.Index:])
		matches := regex.FindStringSubmatch(remaining)

		if matches == nil {
			return "", fmt.Errorf("todo")
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

func Regex(pattern string) Parser {
	return RegexGroup(pattern, 0)
}
