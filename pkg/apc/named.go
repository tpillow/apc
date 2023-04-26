package apc

func Named[CT, T any](name string, parser Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		ctx.PushName(name)
		node, err := parser(ctx)
		ctx.PopName()
		return node, err
	}
}
