package apc

import (
	"fmt"
	"strings"
)

func Eof() Parser {
	return Exact(EofToken{})
}

func Exact(expectToken any) Parser {
	desc := fmt.Sprintf("exactly '%v'", expectToken)
	return NewParser(desc, func(ctx Context) (any, error) {
		token := ctx.Peek()
		if token == expectToken {
			return ctx.Pop(), nil
		}
		return nil, ParseError{
			Up:            nil,
			Expected:      desc,
			Unexpected:    token,
			StartLocation: ctx.CurLocation(),
			EndLocation:   ctx.CurLocation(),
		}
	})
}

func Succeed(value any) Parser {
	return NewParser("always successful", func(ctx Context) (any, error) {
		return value, nil
	})
}

func Fail(err error) Parser {
	return NewParser(fmt.Sprintf("always failing with error '%s'", err), func(ctx Context) (any, error) {
		return nil, err
	})
}

func Seq(parsers ...Parser) Parser {
	if len(parsers) < 1 {
		panic("Seq parser must have at least 1 parser")
	}

	parserDescriptions := []string{}
	for _, parser := range parsers {
		parserDescriptions = append(parserDescriptions, fmt.Sprintf("(%s)", parser.Description))
	}
	desc := fmt.Sprintf("sequence of %d parsers: %s", len(parsers), strings.Join(parserDescriptions, ", "))

	return NewParser(desc, func(ctx Context) (any, error) {
		results := []any{}
		for _, parser := range parsers {
			startLoc := ctx.CurLocation()
			result, err := parser.Parse(ctx)
			if err != nil {
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

func AnyOf(parsers ...Parser) Parser {
	if len(parsers) < 1 {
		panic("Alt parser must have at least 1 parser")
	}

	parserDescriptions := []string{}
	for _, parser := range parsers {
		parserDescriptions = append(parserDescriptions, fmt.Sprintf("(%s)", parser.Description))
	}
	desc := fmt.Sprintf("any of %d parsers: %s", len(parsers), strings.Join(parserDescriptions, ", "))

	return NewParser(desc, func(ctx Context) (any, error) {
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
			Expected:      desc,
			Unexpected:    ctx.Peek(),
			StartLocation: concreteStartLoc,
			EndLocation:   ctx.CurLocation(),
		}
	})
}
