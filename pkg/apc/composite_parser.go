package apc

import "math"

func (parser *Parser) Bind(value any) *Parser {
	return parser.Map(func(_ any) any { return value })
}

func (parser *Parser) Head() *Parser {
	return parser.Index(0)
}

func (parser *Parser) Tail() *Parser {
	return parser.Index(-1)
}

func (first *Parser) Then(second *Parser) *Parser {
	return Seq(first, second).Tail()
}

func (first *Parser) Skip(second *Parser) *Parser {
	return Seq(first, second).Head()
}

func (parser *Parser) AtMost(times int) *Parser {
	return parser.Times(0, times)
}

func (parser *Parser) AtLeast(times int) *Parser {
	return parser.Times(times, math.MaxInt32)
}

func (parser *Parser) Many() *Parser {
	return parser.AtLeast(0)
}

func (first *Parser) Or(second *Parser) *Parser {
	return AnyOf(first, second)
}

func (parser *Parser) SeparatedBy(sepParser *Parser, min int, max int) *Parser {
	if min == 0 {
		return parser.Times(1, 1).ConcatSlices(sepParser.Then(parser).Times(0, max-1)).Or(Succeed([]any{}))
	}
	return parser.Times(1, 1).ConcatSlices(sepParser.Then(parser).Times(min-1, max-1))
}
