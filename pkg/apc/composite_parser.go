package apc

import "math"

func MapValue(parser *Parser, value any) *Parser {
	return Map(parser, func(_ any) any { return value })
}

func MapSlice(parser *Parser, transform MapSliceFunc) *Parser {
	return Map(parser, func(rawValues any) any {
		values, ok := rawValues.([]any)
		if !ok {
			panic("MapSlice parser must have a result of type []any")
		}
		return transform(values)
	})
}

func Head(parser *Parser) *Parser {
	return Index(parser, 0)
}

func Tail(parser *Parser) *Parser {
	return Index(parser, -1)
}

func Then(first *Parser, second *Parser) *Parser {
	return Tail(Seq(first, second))
}

func Skip(first *Parser, second *Parser) *Parser {
	return Head(Seq(first, second))
}

func AtMost(parser *Parser, times int) *Parser {
	return Times(parser, 0, times)
}

func AtLeast(parser *Parser, times int) *Parser {
	return Times(parser, times, math.MaxInt32)
}

func Many(parser *Parser) *Parser {
	return AtLeast(parser, 0)
}

func SeparatedBy(parser *Parser, sepParser *Parser, min int, max int) *Parser {
	if min == 0 {
		return AnyOf(
			ConcatSlices(
				Times(parser, 1, 1),
				Times(Then(sepParser, parser), 0, max-1)),
			Succeed([]any{}))
	}
	return ConcatSlices(Times(parser, 1, 1), Times(Then(sepParser, parser), min-1, max-1))
}
