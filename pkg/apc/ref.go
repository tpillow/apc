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
//	value = OneOf("", MapToAny(Exact("hello")), hashValue)
func Ref[CT, T any](parserPtr *Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		if parserPtr == nil {
			panic("cannot have a Ref to a nil parser")
		}
		return (*parserPtr)(ctx)
	}
}
