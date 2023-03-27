package apc

// MapFunc is a function that maps some type T to some type U. The Origin of node
// is also taken.
type MapFunc[T, U any] func(node T, origin Origin) U

// Returns a parser that maps a Parser[CT, T] into a Parser[CT, U] by running the
// result of parser through mapFunc.
func Map[CT, T, U any](parser Parser[CT, T], mapFunc MapFunc[T, U]) Parser[CT, U] {
	return func(ctx Context[CT]) (U, error) {
		origin := ctx.GetCurOrigin()
		node, err := parser(ctx)
		if err != nil {
			return zeroVal[U](), err
		}
		return mapFunc(node, origin), nil
	}
}

// Returns a parser that maps a Parser[CT, T] into a Parser[CT, U] by always
// returning node.
func Bind[CT, T, U any](parser Parser[CT, T], node U) Parser[CT, U] {
	return Map(parser, func(_ T, _ Origin) U {
		return node
	})
}

// Returns a parser that maps a Parser[CT, T] into a Parser[CT, U] by casting
// the result of parser to type U.
func CastTo[CT, T, U any](parser Parser[CT, T]) Parser[CT, U] {
	return Map(parser, func(node T, _ Origin) U {
		return any(node).(U)
	})
}

// Equivalent to CastTo[CT, T, any](parser).
func CastToAny[CT, T any](parser Parser[CT, T]) Parser[CT, any] {
	return Map(parser, func(node T, _ Origin) any {
		return node
	})
}
