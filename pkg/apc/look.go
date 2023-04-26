package apc

// Provides backtracking support for the provided parser.
// If an error occurs, any consumptions made are reverted to the state
// of the context when this Look parser is called.
// If no error occurs, any consumptions made are committed to the current
// Look frame.
func Look[CT, T any](parser Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		lastLook := ctx.GetLook()
		if lastLook == InvalidLook {
			ctx.SetLook(0)
		}

		node, err := parser(ctx)
		if err != nil {
			org := ctx.GetCurOrigin()
			ctx.SetLook(lastLook)
			return zeroVal[T](), &ParseError{
				Err:     err,
				Message: "",
				Origin:  org,
			}
		}

		newLook := ctx.GetLook()
		if lastLook == InvalidLook {
			ctx.SetLook(InvalidLook)
			_, err := ctx.Consume(newLook)
			if err != nil {
				return zeroVal[T](), err
			}
		} else {
			ctx.SetLook(newLook)
		}

		if err != nil {
			return zeroVal[T](), err
		}
		return node, nil
	}
}
