package apc

import (
	"fmt"
	"regexp"
	"strings"
)

func String(what string) Parser[rune, string] {
	if len(what) <= 0 {
		panic("cannot use String Parser with empty string")
	}
	return NewCallbackParser(func(ctx Context[rune]) (string, error) {
		if ctx.IsEof() {
			return "", ExpErr{
				AtIndex:    ctx.CurIndex(),
				Expected:   fmt.Sprintf("%v", what),
				Unexpected: "EOF",
			}
		}
		remaining := string(ctx.PeekN(ctx.TotalLength() - ctx.CurIndex()))
		if !strings.HasPrefix(remaining, what) {
			return "", ExpErr{
				AtIndex:    ctx.CurIndex(),
				Expected:   what,
				Unexpected: remaining,
			}
		}
		val := string(ctx.PopN(len(what)))
		if val != what {
			panic("unreachable")
		}
		return val, nil
	})
}

func RegexGroup(pattern string, groupIndex int) Parser[rune, string] {
	regex := regexp.MustCompile("^" + pattern)

	return NewCallbackParser(func(ctx Context[rune]) (string, error) {
		if ctx.IsEof() {
			return "", ExpErr{
				AtIndex:    ctx.CurIndex(),
				Expected:   pattern,
				Unexpected: "EOF",
			}
		}
		remaining := string(ctx.PeekN(ctx.TotalLength() - ctx.CurIndex()))
		matches := regex.FindStringSubmatch(remaining)
		if matches == nil {
			return "", ExpErr{
				AtIndex:    ctx.CurIndex(),
				Expected:   fmt.Sprintf("regex<%s>", pattern),
				Unexpected: remaining,
			}
		}
		matchGroup := matches[groupIndex]
		if len(matchGroup) <= 0 {
			panic("a Regex Parser cannot have a 0 length match")
		}
		val := string(ctx.PopN(len(matchGroup)))
		if val != matchGroup {
			panic("unreachable")
		}
		return val, nil
	})
}

func Regex(pattern string) Parser[rune, string] {
	return RegexGroup(pattern, 0)
}
