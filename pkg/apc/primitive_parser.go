package apc

import "fmt"

func Eof() Parser {
	return Exact(EofToken{})
}

func Exact(expectToken any) Parser {
	return NewParser(func(ctx Context) (any, error) {
		token := ctx.Peek()
		if token == expectToken {
			return ctx.Pop(), nil
		}
		return nil, fmt.Errorf("no match")
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
			result, err := parser.Parse(ctx)
			if err != nil {
				return nil, err
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
		for i, parser := range parsers {
			startLoc := ctx.CurLocation()
			result, err := parser.Parse(ctx)
			if err != nil {
				if i == len(parsers)-1 {
					return nil, err
				}
				ctx.SetLocation(startLoc)
				continue
			}
			return result, nil
		}
		return nil, fmt.Errorf("todo")
	})
}
