package apc

func Ref[T any](parserPtr *Parser[T]) Parser[T] {
	return func(ctx Context) (T, error) {
		parser := *parserPtr
		if parser == nil {
			panic("cannot have a Ref to a nil parser")
		}
		return parser(ctx)
	}
}
