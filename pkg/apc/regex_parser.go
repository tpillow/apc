package apc

import (
	"errors"
	"fmt"
	"regexp"
)

// TODO: find a way to not limit max. length of match.
// Might be able to create a RuneReader around ctx.PeekRune?
const regexPeekRuneBufferSize = 1024

type RegexParser struct {
	Regex *regexp.Regexp
}

func Regex(pattern string) *RegexParser {
	if pattern[0] != '^' {
		pattern = fmt.Sprintf("^%v", pattern)
	}
	return &RegexParser{
		Regex: regexp.MustCompile(pattern),
	}
}

func Whitespace() *RegexParser {
	return Regex("\\s+")
}

func (p *RegexParser) Parse(ctx Context) (Node, error) {
	peek, err := PeekNRunes(ctx, 0, regexPeekRuneBufferSize)
	if !errors.Is(err, &EOFError{}) {
		return nil, NewParseError(ctx.GetOrigin(), "%v", err)
	}
	loc := p.Regex.FindStringIndex(peek)
	if loc == nil {
		peekPreview := peek
		if len(peekPreview) > 16 {
			peekPreview = fmt.Sprintf("%v ...more", peekPreview[16:])
		}
		return nil, NewParseError(ctx.GetOrigin(), "expected regex match '%v' but got '%v'", p.Regex, peekPreview)
	}
	if loc[0] != 0 {
		panic("regex pattern must never start match past the start of the string")
	}
	val, err := ConsumeNRunes(ctx, loc[1])
	if err != nil {
		return nil, NewParseError(ctx.GetOrigin(), "expected regex match '%v' but got '%v'", p.Regex, val)
	}
	return val, nil
}

func (p *RegexParser) Map(mapFunc MapFunc) Parser {
	return Map(p, mapFunc)
}
