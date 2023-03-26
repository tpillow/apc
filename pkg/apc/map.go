package apc

// MapFunc is a function that maps some type T to some type U.
type MapFunc[T, U any] func(node T) U

// Returns a parser that maps a Parser[T] into a Parser[U] by running the
// result of the parser through mapFunc.
func Map[CT, T, U any](parser Parser[CT, T], mapFunc MapFunc[T, U]) Parser[CT, U] {
	return func(ctx Context[CT]) (U, error) {
		node, err := parser(ctx)
		if err != nil {
			return zeroVal[U](), err
		}
		return mapFunc(node), nil
	}
}

// Returns a parser that maps Parser[T] to always return node as the result.
func Bind[CT, T, U any](parser Parser[CT, T], node U) Parser[CT, U] {
	return Map(parser, func(_ T) U {
		return node
	})
}

// Returns a parser that maps Parser[T] to return its result casted to type U.
func MapCast[CT, T, U any](parser Parser[CT, T]) Parser[CT, U] {
	return Map(parser, func(node T) U {
		return any(node).(U)
	})
}

// Returns a parser that maps Parser[T] to return its result casted to type any.
func MapToAny[CT, T any](parser Parser[CT, T]) Parser[CT, any] {
	return Map(parser, func(node T) any {
		return node
	})
}
