package apc

// Returns a parser that temporarily adds the skipParser to the Context
// while parsing with parser.
func Skip[T any](skipParser Parser[any], parser Parser[T]) Parser[T] {
	return func(ctx Context) (T, error) {
		ctx.AddSkipParser(skipParser)
		defer ctx.RemoveSkipParser(skipParser)

		node, err := parser(ctx)
		return node, err
	}
}

// Returns a parser that temporarily removes the skipParser from the Context
// while parsing with parser.
func Unskip[T any](skipParser Parser[any], parser Parser[T]) Parser[T] {
	return func(ctx Context) (T, error) {
		ctx.RemoveSkipParser(skipParser)
		defer ctx.AddSkipParser(skipParser)

		node, err := parser(ctx)
		return node, err
	}
}
