package apc

import (
	"fmt"
)

// Returns a parser that runs parser at least min, but at most max, times.
// Returns each parser result in order as a slice.
//
// The min must be >= 0, and max must be > 0. Unless max == -1, in which case
// no maximum is set.
func Range[CT, T any](min int, max int, parser Parser[CT, T]) Parser[CT, []T] {
	if min < 0 {
		panic("min must be >= 0")
	}
	if max != -1 && max <= 0 {
		panic("max must be either -1 (no limit) or > 0")
	}

	return func(ctx Context[CT]) ([]T, error) {
		nodes := make([]T, 0)

		node, err := parser(ctx)
		for err == nil && (max == -1 || len(nodes) < max) {
			nodes = append(nodes, node)
			if max != -1 && len(nodes) >= max {
				break
			}
			node, err = parser(ctx)
		}
		if IsMustReturnParseErr(err) {
			return nil, err
		}

		if len(nodes) < min {
			msg := fmt.Sprintf("at least %v of %v", min, ctx.GetCurParserName())
			if len(nodes) == 0 {
				return nil, ParseErrExpectedButGot(ctx, msg, len(nodes), err)
			}
			return nil, ParseErrConsumedExpectedButGot(ctx, msg, len(nodes), err)
		}
		return nodes, nil
	}
}

// Same as Range(0, -1, parser).
func ZeroOrMore[CT, T any](parser Parser[CT, T]) Parser[CT, []T] {
	return Range(0, -1, parser)
}

// Same as Range(1, -1, parser).
func OneOrMore[CT, T any](parser Parser[CT, T]) Parser[CT, []T] {
	return Range(1, -1, parser)
}

// Same as Range(0, 1, parser), but with the resulting slice mapped
// to a single value, or default T if 0 matches occurred.
func Maybe[CT, T any](parser Parser[CT, T]) Parser[CT, T] {
	return Map(Range(0, 1, parser), func(node []T) T {
		if node == nil || len(node) <= 0 {
			return zeroVal[T]()
		}
		if len(node) != 1 {
			panic("unreachable: Range(0, 1) should return at most 1 node")
		}
		return node[0]
	})
}

// Same as OneOrMore(parser), but ensures that each subsequent match is separated by
// a successful parse by sepParser. The results of sepParser are not returned.
func OneOrMoreSeparated[CT, T, U any](parser Parser[CT, T],
	sepParser Parser[CT, U]) Parser[CT, []T] {
	sepParse := Map(
		Seq2(sepParser, parser),
		func(node *Seq2Node[U, T]) T {
			return node.Result2
		})

	return Map(
		Seq2(parser, ZeroOrMore(sepParse)),
		func(node *Seq2Node[T, []T]) []T {
			result := []T{node.Result1}
			return append(result, node.Result2...)
		})
}

// Same as ZeroOrMore(parser), but ensures that each subsequent match is separated by
// a successful parse by sepParser. The results of sepParser are not returned.
func ZeroOrMoreSeparated[CT, T, U any](parser Parser[CT, T],
	sepParser Parser[CT, U]) Parser[CT, []T] {
	return Map(
		Maybe(OneOrMoreSeparated(parser, sepParser)),
		func(node []T) []T {
			if node == nil {
				return []T{}
			}
			return node
		})
}
