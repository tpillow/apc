package apc

type MapFunc[T, U any] func(node T) U

type MapParser[T, U any] struct {
	Parser  Parser[T]
	MapFunc MapFunc[T, U]
}

func Map[T, U any](parser Parser[T], mapFunc MapFunc[T, U]) Parser[U] {
	return func(ctx Context) (U, error) {
		node, err := parser(ctx)
		if err != nil {
			return zeroVal[U](), err
		}
		return mapFunc(node), nil
	}
}

func Bind[T, U any](parser Parser[T], node U) Parser[U] {
	return Map(parser, func(_ T) U {
		return node
	})
}

func MapCast[T, U any](parser Parser[T]) Parser[U] {
	return Map(parser, func(node T) U {
		return any(node).(U)
	})
}

func MapToAny[T any](parser Parser[T]) Parser[any] {
	return Map(parser, func(node T) any {
		return node
	})
}
