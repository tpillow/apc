package apc

// MapFunc is a function that maps some type T to some type U.
type MapFunc[T, U any] func(node T) U

// Returns a parser that maps a Parser[T] into a Parser[U] by running the
// result of the parser through mapFunc.
func Map[T, U any](parser Parser[T], mapFunc MapFunc[T, U]) Parser[U] {
	return func(ctx Context) (U, error) {
		node, err := parser(ctx)
		if err != nil {
			return zeroVal[U](), err
		}
		return mapFunc(node), nil
	}
}

// Returns a parser that maps Parser[T] to always return node as the result.
func Bind[T, U any](parser Parser[T], node U) Parser[U] {
	return Map(parser, func(_ T) U {
		return node
	})
}

// Returns a parser that maps Parser[T] to return its result casted to type U.
func MapCast[T, U any](parser Parser[T]) Parser[U] {
	return Map(parser, func(node T) U {
		return any(node).(U)
	})
}

// Returns a parser that maps Parser[T] to return its result casted to type any.
func MapToAny[T any](parser Parser[T]) Parser[any] {
	return Map(parser, func(node T) any {
		return node
	})
}
