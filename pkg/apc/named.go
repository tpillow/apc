package apc

// Associates a name with the provided parser for better error messages.
func Named[CT, T any](name string, parser Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		lastName := ctx.GetCurName()
		ctx.SetCurName(name)
		node, err := parser(ctx)
		ctx.SetCurName(lastName)
		return node, err
	}
}
