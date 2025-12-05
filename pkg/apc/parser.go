package apc

import (
	"fmt"
)

type ParseFunc func(ctx Context) (any, error)

type Parser struct {
	Description string
	ParseFunc   ParseFunc
}

func NewParser(description string, parseFunc ParseFunc) *Parser {
	return &Parser{
		Description: description,
		ParseFunc:   parseFunc,
	}
}

func NewFutureParser() *Parser {
	return NewParser("", nil)
}

func (parser *Parser) Parse(ctx Context) (any, error) {
	if parser.ParseFunc == nil {
		panic("cannot use a Parser whose ParseFunc is nil")
	}
	return parser.ParseFunc(ctx)
}

func (parser *Parser) ParseToEof(ctx Context) (any, error) {
	return parser.Skip(Eof).Parse(ctx)
}

func (parser *Parser) Peek() *Parser {
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

func (parser *Parser) Map(transform func(value any) any) *Parser {
	return NewParser(parser.Description, func(ctx Context) (any, error) {
		result, err := parser.Parse(ctx)
		if err != nil {
			return nil, err
		}
		return transform(result), nil
	})
}

func (parser *Parser) Index(index int) *Parser {
	return parser.MapSlice(func(values []any) any {
		if index < 0 {
			index = len(values) + index
		}
		if index < 0 {
			panic("cannot use Index with out of bounds index < 0")
		}
		if index >= len(values) {
			panic("cannot use Index with out of bounds index >= len([]any)")
		}
		return values[index]
	})
}

// TODO: rethink this? genericize? multi-param?
func (first *Parser) ConcatSlices(second *Parser) *Parser {
	return Seq(first, second).MapSlice(func(results []any) any {
		resultA, ok := results[0].([]any)
		if !ok {
			panic("ConcatSlices input parsers must produce a result of type []any")
		}
		resultB, ok := results[1].([]any)
		if !ok {
			panic("ConcatSlices input parsers must produce a result of type []any")
		}
		return append(resultA, resultB...)
	})
}

func (parser *Parser) Optional(defaultValue any) *Parser {
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

func (parser *Parser) Describe(description string) *Parser {
	parser.Description = description
	return parser
}

func (parser *Parser) Become(other *Parser) {
	if parser.ParseFunc != nil {
		panic("cannot call Become on a parser with a non-nil ParseFunc")
	}
	parser.Description = other.Description
	parser.ParseFunc = other.ParseFunc
}

func (parser *Parser) Times(min int, max int) *Parser {
	if min < 0 {
		panic("Times parser min must be >= 0")
	}
	if max < 0 {
		panic("Times parser max must be >= 0")
	}
	if max < min {
		panic("Times parser max must be >= min")
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
