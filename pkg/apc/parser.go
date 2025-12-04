package apc

import (
	"fmt"
	"math"
)

type ParseFunc func(ctx Context) (any, error)

type Parser struct {
	Description string
	ParseFunc   ParseFunc
}

func NewParser(description string, parseFunc ParseFunc) Parser {
	return Parser{
		Description: description,
		ParseFunc:   parseFunc,
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
	desc := fmt.Sprintf("peeking parser of %s", parser.Description)
	return NewParser(desc, func(ctx Context) (any, error) {
		startLoc := ctx.CurLocation()
		result, err := parser.Parse(ctx)
		endLoc := ctx.CurLocation()
		ctx.SetLocation(startLoc)
		return result, ParseError{
			Up:            err,
			Expected:      desc,
			Unexpected:    ctx.Peek(),
			StartLocation: startLoc,
			EndLocation:   endLoc,
		}
	})
}

func (parser Parser) Map(transform func(value any) any) Parser {
	return NewParser(parser.Description, func(ctx Context) (any, error) {
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
	desc := fmt.Sprintf("optional %s", parser.Description)
	return NewParser(desc, func(ctx Context) (any, error) {
		startLoc := ctx.CurLocation()
		result, err := parser.Parse(ctx)
		if err != nil {
			ctx.SetLocation(startLoc)
			return defaultValue, nil
		}
		return result, nil
	})
}

func (parser Parser) Describe(description string) Parser {
	return NewParser(description, func(ctx Context) (any, error) {
		startLoc := ctx.CurLocation()
		result, err := parser.Parse(ctx)
		if err != nil {
			return nil, ParseError{
				Up:            err,
				Expected:      description,
				Unexpected:    ctx.Peek(),
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}
		return result, nil
	})
}

func (first Parser) Then(second Parser) Parser {
	desc := fmt.Sprintf("%s followed by %s, keeping the second", first.Description, second.Description)

	return NewParser(desc, func(ctx Context) (any, error) {
		startLoc := ctx.CurLocation()
		_, err := first.Parse(ctx)
		if err != nil {
			return nil, ParseError{
				Up:            err,
				Expected:      desc,
				Unexpected:    ctx.Peek(),
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}

		startLoc = ctx.CurLocation()
		result, err := second.Parse(ctx)
		if err != nil {
			return nil, ParseError{
				Up:            err,
				Expected:      desc,
				Unexpected:    ctx.Peek(),
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}
		return result, nil
	})
}

func (first Parser) Skip(second Parser) Parser {
	desc := fmt.Sprintf("%s followed by %s, keeping the first", first.Description, second.Description)
	return NewParser(desc, func(ctx Context) (any, error) {
		startLoc := ctx.CurLocation()
		firstResult, err := first.Parse(ctx)
		if err != nil {
			return nil, ParseError{
				Up:            err,
				Expected:      desc,
				Unexpected:    ctx.Peek(),
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}

		startLoc = ctx.CurLocation()
		_, err = second.Parse(ctx)
		if err != nil {
			return nil, ParseError{
				Up:            err,
				Expected:      desc,
				Unexpected:    ctx.Peek(),
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
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

	desc := fmt.Sprintf("%s %d to %d times", parser.Description, min, max)

	return NewParser(desc, func(ctx Context) (any, error) {
		results := []any{}

		for len(results) < max {
			startLoc := ctx.CurLocation()
			result, err := parser.Parse(ctx)

			if err != nil {
				if len(results) >= min {
					ctx.SetLocation(startLoc)
					break
				}
				return nil, ParseError{
					Up:            err,
					Expected:      desc,
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

func (parser Parser) AtMost(times int) Parser {
	return parser.Times(0, times)
}

func (parser Parser) AtLeast(times int) Parser {
	return parser.Times(times, math.MaxInt32)
}

func (parser Parser) Many() Parser {
	return parser.AtLeast(0)
}
