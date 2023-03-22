package apc

import "errors"

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

func seqSetResultHelper[T any](first bool, ctx Context, name string, parser Parser[T], resultField *T) error {
	node, err := parser(ctx)
	if err != nil {
		if !errors.Is(err, ErrParseErr) {
			return err
		}
		if first {
			return ParseErrExpectedButGotNext(ctx, name, err)
		}
		return ParseErrConsumedExpectedButGotNext(ctx, name, err)
	}
	*resultField = node
	return nil
}

type Seq2Node[T1, T2 any] struct {
	Result1 T1
	Result2 T2
}

func Seq2[T1, T2 any](name string, parser1 Parser[T1], parser2 Parser[T2]) Parser[Seq2Node[T1, T2]] {
	return func(ctx Context) (Seq2Node[T1, T2], error) {
		result := Seq2Node[T1, T2]{}

		if err := seqSetResultHelper(true, ctx, name, parser1, &result.Result1); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser2, &result.Result2); err != nil {
			return result, err
		}

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

		if err := seqSetResultHelper(true, ctx, name, parser1, &result.Result1); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser2, &result.Result2); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser3, &result.Result3); err != nil {
			return result, err
		}

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

		if err := seqSetResultHelper(true, ctx, name, parser1, &result.Result1); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser2, &result.Result2); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser3, &result.Result3); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser4, &result.Result4); err != nil {
			return result, err
		}

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

		if err := seqSetResultHelper(true, ctx, name, parser1, &result.Result1); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser2, &result.Result2); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser3, &result.Result3); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser4, &result.Result4); err != nil {
			return result, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser5, &result.Result5); err != nil {
			return result, err
		}

		return result, nil
	}
}