package apc

func Skip[T any](skipParser Parser[any], parser Parser[T]) Parser[T] {
	return func(ctx Context) (T, error) {
		ctx.AddSkipParser(skipParser)
		defer ctx.RemoveSkipParser(skipParser)
		ctx.RunSkipParsers()

		node, err := parser(ctx)
		return node, err
	}
}

func Unskip[T any](skipParser Parser[any], parser Parser[T]) Parser[T] {
	return func(ctx Context) (T, error) {
		ctx.RemoveSkipParser(skipParser)
		defer ctx.RunSkipParsers()
		defer ctx.AddSkipParser(skipParser)

		node, err := parser(ctx)
		return node, err
	}
}
