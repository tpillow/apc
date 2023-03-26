package apc

// Returns a parser that parses all provided parsers in order.
// Returns each parser result in order as a slice.
//
// The number of parsers provided must be at least 2.
func Seq[CT, T any](name string, parsers ...Parser[CT, T]) Parser[CT, []T] {
	if len(parsers) < 2 {
		panic("must provide at least 2 parsers to Seq")
	}

	return func(ctx Context[CT]) ([]T, error) {
		nodes := make([]T, 0)
		for i, parser := range parsers {
			node, err := parser(ctx)
			if err != nil {
				if IsMustReturnParseErr(err) {
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

// Internal helper function used with Seq# parsers.
func seqSetResultHelper[CT, T any](first bool, ctx Context[CT], name string, parser Parser[CT, T], resultField *T) error {
	node, err := parser(ctx)
	if err != nil {
		if IsMustReturnParseErr(err) {
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

// Seq2Node holds 2 generically-typed results.
type Seq2Node[T1, T2 any] struct {
	Result1 T1
	Result2 T2
}

// Returns a parser that parses all provided parsers in order.
// This is the same as Seq, but is optimized for N parsers of different types.
// Returns each parser result in the corresponding typed result field.
func Seq2[CT, T1, T2 any](name string, parser1 Parser[CT, T1], parser2 Parser[CT, T2]) Parser[CT, *Seq2Node[T1, T2]] {
	return func(ctx Context[CT]) (*Seq2Node[T1, T2], error) {
		result := &Seq2Node[T1, T2]{}

		if err := seqSetResultHelper(true, ctx, name, parser1, &result.Result1); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser2, &result.Result2); err != nil {
			return nil, err
		}

		return result, nil
	}
}

// Seq3Node holds 3 generically-typed results.
type Seq3Node[T1, T2, T3 any] struct {
	Result1 T1
	Result2 T2
	Result3 T3
}

// Returns a parser that parses all provided parsers in order.
// This is the same as Seq, but is optimized for N parsers of different types.
// Returns each parser result in the corresponding typed result field.
func Seq3[CT, T1, T2, T3 any](name string, parser1 Parser[CT, T1], parser2 Parser[CT, T2],
	parser3 Parser[CT, T3]) Parser[CT, *Seq3Node[T1, T2, T3]] {

	return func(ctx Context[CT]) (*Seq3Node[T1, T2, T3], error) {
		result := &Seq3Node[T1, T2, T3]{}

		if err := seqSetResultHelper(true, ctx, name, parser1, &result.Result1); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser2, &result.Result2); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser3, &result.Result3); err != nil {
			return nil, err
		}

		return result, nil
	}
}

// Seq4Node holds 4 generically-typed results.
type Seq4Node[T1, T2, T3, T4 any] struct {
	Result1 T1
	Result2 T2
	Result3 T3
	Result4 T4
}

// Returns a parser that parses all provided parsers in order.
// This is the same as Seq, but is optimized for N parsers of different types.
// Returns each parser result in the corresponding typed result field.
func Seq4[CT, T1, T2, T3, T4 any](name string, parser1 Parser[CT, T1], parser2 Parser[CT, T2],
	parser3 Parser[CT, T3], parser4 Parser[CT, T4]) Parser[CT, *Seq4Node[T1, T2, T3, T4]] {

	return func(ctx Context[CT]) (*Seq4Node[T1, T2, T3, T4], error) {
		result := &Seq4Node[T1, T2, T3, T4]{}

		if err := seqSetResultHelper(true, ctx, name, parser1, &result.Result1); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser2, &result.Result2); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser3, &result.Result3); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser4, &result.Result4); err != nil {
			return nil, err
		}

		return result, nil
	}
}

// Seq5Node holds 5 generically-typed results.
type Seq5Node[T1, T2, T3, T4, T5 any] struct {
	Result1 T1
	Result2 T2
	Result3 T3
	Result4 T4
	Result5 T5
}

// Returns a parser that parses all provided parsers in order.
// This is the same as Seq, but is optimized for N parsers of different types.
// Returns each parser result in the corresponding typed result field.
func Seq5[CT, T1, T2, T3, T4, T5 any](name string, parser1 Parser[CT, T1], parser2 Parser[CT, T2],
	parser3 Parser[CT, T3], parser4 Parser[CT, T4], parser5 Parser[CT, T5]) Parser[CT, *Seq5Node[T1, T2, T3, T4, T5]] {

	return func(ctx Context[CT]) (*Seq5Node[T1, T2, T3, T4, T5], error) {
		result := &Seq5Node[T1, T2, T3, T4, T5]{}

		if err := seqSetResultHelper(true, ctx, name, parser1, &result.Result1); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser2, &result.Result2); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser3, &result.Result3); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser4, &result.Result4); err != nil {
			return nil, err
		}
		if err := seqSetResultHelper(false, ctx, name, parser5, &result.Result5); err != nil {
			return nil, err
		}

		return result, nil
	}
}
