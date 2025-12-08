package apc

import "math"

func MapValue[IT, OT, OT2 any](parser *Parser[IT, OT], value OT2) *Parser[IT, OT2] {
	return Map(parser, func(_ OT) OT2 { return value })
}

func MapSlice[IT, OT, OT2 any](parser *Parser[IT, []OT], transform MapSliceFunc[OT, OT2]) *Parser[IT, OT2] {
	return Map(parser, func(values []OT) OT2 {
		return transform(values)
	})
}

func Head[IT, OT any](parser *Parser[IT, []OT]) *Parser[IT, OT] {
	return Index(parser, 0)
}

func Tail[IT, OT any](parser *Parser[IT, []OT]) *Parser[IT, OT] {
	return Index(parser, -1)
}

func Then[IT, OT, OT2 any](first *Parser[IT, OT], second *Parser[IT, OT2]) *Parser[IT, OT2] {
	return Map(Seq2(first, second), func(results Result2[OT, OT2]) OT2 {
		return results.Value2
	})
}

func Skip[IT, OT, OT2 any](first *Parser[IT, OT], second *Parser[IT, OT2]) *Parser[IT, OT] {
	return Map(Seq2(first, second), func(results Result2[OT, OT2]) OT {
		return results.Value1
	})
}

func AtMost[IT, OT any](parser *Parser[IT, OT], times int) *Parser[IT, []OT] {
	return Times(parser, 0, times)
}

func AtLeast[IT, OT any](parser *Parser[IT, OT], times int) *Parser[IT, []OT] {
	return Times(parser, times, math.MaxInt32)
}

func Many[IT, OT any](parser *Parser[IT, OT]) *Parser[IT, []OT] {
	return AtLeast(parser, 0)
}

func SeparatedBy[IT, OT, OT2 any](parser *Parser[IT, OT], sepParser *Parser[IT, OT2], min int, max int) *Parser[IT, []OT] {
	if min == 0 {
		return AnyOf(
			ConcatSlices(
				Times(parser, 1, 1),
				Times(Then(sepParser, parser), 0, max-1)),
			Succeed[IT]([]OT{}))
	}
	return ConcatSlices(Times(parser, 1, 1), Times(Then(sepParser, parser), min-1, max-1))
}
