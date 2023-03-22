package apc

import (
	"errors"
	"fmt"
)

func Range[T any](name string, min int, max int, parser Parser[T]) Parser[[]T] {
	return func(ctx Context) ([]T, error) {
		nodes := make([]T, 0)

		node, err := parser(ctx)
		for err == nil {
			nodes = append(nodes, node)
			node, err = parser(ctx)
		}
		if !errors.Is(err, ErrParseErr) {
			return nil, err
		}

		if len(nodes) < min {
			return nil, ParseErrExpectedButGot(ctx, fmt.Sprintf("at least %v of %v", min, name), len(nodes), err)
		}
		if max != -1 && len(nodes) > max {
			return nil, ParseErrExpectedButGot(ctx, fmt.Sprintf("at most %v of %v", max, name), len(nodes), err)
		}
		return nodes, nil
	}
}

func ZeroOrMore[T any](name string, parser Parser[T]) Parser[[]T] {
	return Range(name, 0, -1, parser)
}

func OneOrMore[T any](name string, parser Parser[T]) Parser[[]T] {
	return Range(name, 1, -1, parser)
}

func Maybe[T any](name string, parser Parser[T]) Parser[T] {
	return Map(Range(name, 0, 1, parser), func(node []T) T {
		if node == nil || len(node) <= 0 {
			return zeroVal[T]()
		}
		if len(node) != 1 {
			panic("unreachable: Range(0, 1) should return at most 1 node")
		}
		return node[0]
	})
}

func OneOrMoreSeparated[T, U any](name string, parser Parser[T], sepParser Parser[U]) Parser[[]T] {
	sepParse := Map(
		Seq2(name, sepParser, parser),
		func(node Seq2Node[U, T]) T {
			return node.Result2
		})

	return Map(
		Seq2(name, parser, ZeroOrMore("", sepParse)),
		func(node Seq2Node[T, []T]) []T {
			result := []T{node.Result1}
			return append(result, node.Result2...)
		})
}

func ZeroOrMoreSeparated[T, U any](name string, parser Parser[T], sepParser Parser[U]) Parser[[]T] {
	return Map(
		Maybe(name, OneOrMoreSeparated(name, parser, sepParser)),
		func(node []T) []T {
			if node == nil {
				return []T{}
			}
			return node
		})
}
