package apc

// Provides backtracking support for the provided parser.
// If an error occurs, any consumptions made are reverted to the state
// of the context when this Look parser is called.
// If no error occurs, any consumptions made are committed to the current
// Look frame.
func Look[CT, T any](parser Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		lastLook := ctx.GetLookOffset()
		if lastLook == InvalidLookOffset {
			ctx.SetLookOffset(0)
		}

		node, err := parser(ctx)
		if err != nil {
			org := ctx.GetCurOrigin()
			ctx.SetLookOffset(lastLook)
			if pec, ok := err.(*ParseErrorConsumed); ok {
				// Just convert the ParseErrorConsumed to a ParseError.
				return zeroVal[T](), &ParseError{
					Err:     pec.Err,
					Message: pec.Message,
					Origin:  pec.Origin,
				}
			}
			return zeroVal[T](), &ParseError{
				Err:     err,
				Message: "",
				Origin:  org,
			}
		}

		newLook := ctx.GetLookOffset()
		if lastLook == InvalidLookOffset {
			ctx.SetLookOffset(InvalidLookOffset)
			_, err := ctx.Consume(newLook)
			if err != nil {
				return zeroVal[T](), err
			}
		} else {
			ctx.SetLookOffset(newLook)
		}

		return node, err
	}
}
