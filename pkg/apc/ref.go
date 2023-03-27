package apc

// Creates a Parser[CT, T] from a *Parser[CT, T].
//
// Useful for avoiding circular variable dependencies. For example:
//
//	var value = OneOf("", MapToAny(ExactStr("hello")), MapToAny(hashValue))
//	var hashValue = Seq("", MapToAny(ExactStr("#")), value)
//
// Is invalid, however this can be remedied by:
//
//	var value Parser[rune, any]
//	var valueRef = Ref(&value)
//	var hashValue = Seq("", MapToAny(ExactStr("#")), valueRef)
//
//	// At runtime, in some initialization function:
//	value = OneOf("", MapToAny(ExactStr("hello")), hashValue)
func Ref[CT, T any](parserPtr *Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		if parserPtr == nil {
			panic("cannot have a Ref to a nil parser")
		}
		return (*parserPtr)(ctx)
	}
}
