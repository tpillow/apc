package apc

func Look[CT, T any](parser Parser[CT, T]) Parser[CT, T] {
	return func(ctx Context[CT]) (T, error) {
		lookCtx, ok := ctx.(LookContext)
		if !ok {
			panic("cannot Look() with a Context that is not a LookContext")
		}
		lookCtx.NewLook()
		node, err := parser(ctx)
		if err != nil {
			org := ctx.GetCurOrigin()
			lookCtx.RevertLook()
			return zeroVal[T](), &ParseError{
				Err:     err,
				Message: "",
				Origin:  org,
			}
		}
		err = lookCtx.CommitLook()
		if err != nil {
			return zeroVal[T](), err
		}
		return node, nil
	}
}
