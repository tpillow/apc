package apc

import "errors"

type SeqParser[T any] struct {
	Name    string
	Parsers []Parser[T]
}

func Seq[T any](name string, parsers ...Parser[T]) Parser[[]T] {
	if len(parsers) < 2 {
		panic("must provide at least 2 parsers to Seq")
	}

	return func(ctx Context) ([]T, error) {
		nodes := make([]T, 0)
		for i, parser := range parsers {
			node, err := parser(ctx)
			if err != nil {
				if !errors.Is(err, ErrParseErr) {
					return nil, err
				}
				if i == 0 {
					return nil, ParseErrExpectedButGotNext(ctx, name, err)
				}
				return nil, ParseErrConsumedExpectedButGotNext(ctx, name, err)
			}
			nodes = append(nodes, node)
		}
		return nodes, nil
	}
}

type Seq2Node[T1, T2 any] struct {
	Result1 T1
	Result2 T2
}

func Seq2[T1, T2 any](name string, parser1 Parser[T1], parser2 Parser[T2]) Parser[Seq2Node[T1, T2]] {
	return func(ctx Context) (Seq2Node[T1, T2], error) {
		result := Seq2Node[T1, T2]{}

		node1, err := parser1(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrExpectedButGotNext(ctx, name, err)
		}
		result.Result1 = node1

		node2, err := parser2(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result2 = node2

		return result, nil
	}
}

type Seq3Node[T1, T2, T3 any] struct {
	Result1 T1
	Result2 T2
	Result3 T3
}

func Seq3[T1, T2, T3 any](name string, parser1 Parser[T1], parser2 Parser[T2],
	parser3 Parser[T3]) Parser[Seq3Node[T1, T2, T3]] {

	return func(ctx Context) (Seq3Node[T1, T2, T3], error) {
		result := Seq3Node[T1, T2, T3]{}

		node1, err := parser1(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrExpectedButGotNext(ctx, name, err)
		}
		result.Result1 = node1

		node2, err := parser2(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result2 = node2

		node3, err := parser3(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result3 = node3

		return result, nil
	}
}

type Seq4Node[T1, T2, T3, T4 any] struct {
	Result1 T1
	Result2 T2
	Result3 T3
	Result4 T4
}

func Seq4[T1, T2, T3, T4 any](name string, parser1 Parser[T1], parser2 Parser[T2],
	parser3 Parser[T3], parser4 Parser[T4]) Parser[Seq4Node[T1, T2, T3, T4]] {

	return func(ctx Context) (Seq4Node[T1, T2, T3, T4], error) {
		result := Seq4Node[T1, T2, T3, T4]{}

		node1, err := parser1(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrExpectedButGotNext(ctx, name, err)
		}
		result.Result1 = node1

		node2, err := parser2(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result2 = node2

		node3, err := parser3(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result3 = node3

		node4, err := parser4(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result4 = node4

		return result, nil
	}
}

type Seq5Node[T1, T2, T3, T4, T5 any] struct {
	Result1 T1
	Result2 T2
	Result3 T3
	Result4 T4
	Result5 T5
}

func Seq5[T1, T2, T3, T4, T5 any](name string, parser1 Parser[T1], parser2 Parser[T2],
	parser3 Parser[T3], parser4 Parser[T4], parser5 Parser[T5]) Parser[Seq5Node[T1, T2, T3, T4, T5]] {

	return func(ctx Context) (Seq5Node[T1, T2, T3, T4, T5], error) {
		result := Seq5Node[T1, T2, T3, T4, T5]{}

		node1, err := parser1(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrExpectedButGotNext(ctx, name, err)
		}
		result.Result1 = node1

		node2, err := parser2(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result2 = node2

		node3, err := parser3(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result3 = node3

		node4, err := parser4(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result4 = node4

		node5, err := parser5(ctx)
		if err != nil {
			if !errors.Is(err, ErrParseErr) {
				return result, err
			}
			return result, ParseErrConsumedExpectedButGotNext(ctx, name, err)
		}
		result.Result5 = node5

		return result, nil
	}
}
