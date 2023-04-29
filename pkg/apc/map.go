package apc

// MapFunc is a function that maps some type T to some type U.
type MapFunc[T, U any] func(node T) U

// Returns a parser that maps a Parser[CT, T] into a Parser[CT, U] by running the
// result of parser through mapFunc.
func Map[CT, T, U any](parser Parser[CT, T], mapFunc MapFunc[T, U]) Parser[CT, U] {
	return func(ctx Context[CT]) (U, error) {
		node, err := parser(ctx)
		if err != nil {
			return zeroVal[U](), err
		}
		return mapFunc(node), nil
	}
}

// Returns a parser that maps a Parser[CT, T] into a Parser[CT, U] by always
// returning node.
func Bind[CT, T, U any](parser Parser[CT, T], node U) Parser[CT, U] {
	return Map(parser, func(_ T) U {
		return node
	})
}

// Returns a parser that maps a Parser[CT, T] into a Parser[CT, U] by casting
// the result of parser to type U.
func CastTo[CT, T, U any](parser Parser[CT, T]) Parser[CT, U] {
	return Map(parser, func(node T) U {
		return any(node).(U)
	})
}

// Equivalent to CastTo[CT, T, any](parser).
func CastToAny[CT, T any](parser Parser[CT, T]) Parser[CT, any] {
	return Map(parser, func(node T) any {
		return node
	})
}
