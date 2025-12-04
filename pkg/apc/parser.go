package apc

import "math"

type ParseFunc func(ctx Context) (any, error)

type Parser struct {
	ParseFunc ParseFunc
}

func NewParser(parseFunc ParseFunc) Parser {
	return Parser{
		ParseFunc: parseFunc,
	}
}

func (parser Parser) Parse(ctx Context) (any, error) {
	if parser.ParseFunc == nil {
		panic("cannot use a Parser whose ParseFunc is nil")
	}
	return parser.ParseFunc(ctx)
}

func (parser Parser) ParseToEof(ctx Context) (any, error) {
	return parser.Skip(Eof()).Parse(ctx)
}

func (parser Parser) Peek() Parser {
	return NewParser(func(ctx Context) (any, error) {
		startLoc := ctx.CurLocation()
		result, err := parser.Parse(ctx)
		ctx.SetLocation(startLoc)
		return result, err
	})
}

func (parser Parser) Map(transform func(value any) any) Parser {
	return NewParser(func(ctx Context) (any, error) {
		result, err := parser.Parse(ctx)
		if err != nil {
			return nil, err
		}
		return transform(result), nil
	})
}

func (parser Parser) Bind(value any) Parser {
	return parser.Map(func(_ any) any { return value })
}

func (parser Parser) Optional(defaultValue any) Parser {
	return NewParser(func(ctx Context) (any, error) {
		startLoc := ctx.CurLocation()
		result, err := parser.Parse(ctx)
		if err != nil {
			ctx.SetLocation(startLoc)
			return defaultValue, nil
		}
		return result, nil
	})
}

func (first Parser) Then(second Parser) Parser {
	return NewParser(func(ctx Context) (any, error) {
		_, err := first.Parse(ctx)
		if err != nil {
			return nil, err
		}
		return second.Parse(ctx)
	})
}

func (first Parser) Skip(second Parser) Parser {
	return NewParser(func(ctx Context) (any, error) {
		firstResult, err := first.Parse(ctx)
		if err != nil {
			return nil, err
		}
		_, err = second.Parse(ctx)
		if err != nil {
			return nil, err
		}
		return firstResult, nil
	})
}

func (parser Parser) Times(min int, max int) Parser {
	if min < 0 || max < 0 {
		panic("Times parser min and max must be >= 0")
	}
	if max == 0 {
		panic("Times parser max must be > 0")
	}
	if max < min {
		panic("Times parser max must be < min")
	}

	return NewParser(func(ctx Context) (any, error) {
		results := []any{}

		for len(results) < max {
			startLoc := ctx.CurLocation()
			result, err := parser.Parse(ctx)

			if err != nil {
				if len(results) >= min {
					ctx.SetLocation(startLoc)
					break
				}
				return nil, err
			}
			results = append(results, result)
		}

		return results, nil
	})
}

func (parser Parser) AtMost(times int) Parser {
	return parser.Times(0, times)
}

func (parser Parser) AtLeast(times int) Parser {
	return parser.Times(times, math.MaxInt32)
}

func (parser Parser) Many() Parser {
	return parser.AtLeast(0)
}
