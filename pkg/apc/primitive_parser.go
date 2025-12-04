package apc

func Eof() Parser {
	return Exact(EofToken{})
}

func Exact(expectToken any) Parser {
	return NewParser(func(ctx Context) (any, error) {
		token := ctx.Peek()
		if token == expectToken {
			return ctx.Pop(), nil
		}
		return nil, ParseError{
			Up:            nil,
			Expected:      expectToken,
			Unexpected:    token,
			StartLocation: ctx.CurLocation(),
			EndLocation:   ctx.CurLocation(),
		}
	})
}

func Succeed(value any) Parser {
	return NewParser(func(ctx Context) (any, error) {
		return value, nil
	})
}

func Fail(err error) Parser {
	return NewParser(func(ctx Context) (any, error) {
		return nil, err
	})
}

func Seq(parsers ...Parser) Parser {
	if len(parsers) < 1 {
		panic("Seq parser must have at least 1 parser")
	}

	return NewParser(func(ctx Context) (any, error) {
		results := []any{}
		for _, parser := range parsers {
			startLoc := ctx.CurLocation()
			result, err := parser.Parse(ctx)
			if err != nil {
				return nil, ParseError{
					Up:            err,
					Expected:      "sequence of parsers to succeed",
					Unexpected:    ctx.Peek(),
					StartLocation: startLoc,
					EndLocation:   ctx.CurLocation(),
				}
			}
			results = append(results, result)
		}
		return results, nil
	})
}

func AnyOf(parsers ...Parser) Parser {
	if len(parsers) < 1 {
		panic("Alt parser must have at least 1 parser")
	}

	return NewParser(func(ctx Context) (any, error) {
		var err error
		concreteStartLoc := ctx.CurLocation()

		for i, parser := range parsers {
			startLoc := ctx.CurLocation()
			var result any
			result, err = parser.Parse(ctx)
			if err != nil {
				if i == len(parsers)-1 {
					break
				}
				ctx.SetLocation(startLoc)
				continue
			}
			return result, nil
		}

		return nil, ParseError{
			Up:            err,
			Expected:      "any one parser to succeed (error shown is last attempted error)",
			Unexpected:    ctx.Peek(),
			StartLocation: concreteStartLoc,
			EndLocation:   ctx.CurLocation(),
		}
	})
}
