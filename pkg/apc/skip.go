package apc

// Returns a parser that temporarily adds the skipParser to the Context
// while parsing with parser.
func Skip[CT any, T any](skipParser Parser[CT, any], parser Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		ctx.AddSkipParser(skipParser)
		defer ctx.RemoveSkipParser(skipParser)

		node, err := parser(ctx)
		return node, err
	}
}

// Returns a parser that temporarily removes the skipParser from the Context
// while parsing with parser.
func Unskip[CT any, T any](skipParser Parser[CT, any], parser Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		ctx.RemoveSkipParser(skipParser)
		defer ctx.AddSkipParser(skipParser)

		node, err := parser(ctx)
		return node, err
	}
}
