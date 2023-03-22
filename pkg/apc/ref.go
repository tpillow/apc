package apc

// Creates a Parser[T] from a *Parser[T].
//
// Useful for avoiding circular variable dependencies. For example:
//
//	var value = OneOf("", MapToAny(Exact("hello")), MapToAny(hashValue))
//	var hashValue = Seq("", MapToAny(Exact("#")), value)
//
// Is invalid, however this can be remedied by:
//
//	var value Parser[any]
//	var valueRef = Ref(&value)
//	var hashValue = Seq("", MapToAny(Exact("#")), valueRef)
//
//	// At runtime, in some initialization function:
//	value = OneOf("", MapToAny(Exact("hello")), MapToAny(hashValue))
func Ref[T any](parserPtr *Parser[T]) Parser[T] {
	return func(ctx Context) (T, error) {
		parser := *parserPtr
		if parser == nil {
			panic("cannot have a Ref to a nil parser")
		}
		return parser(ctx)
	}
}
