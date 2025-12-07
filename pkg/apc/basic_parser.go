package apc

import (
	"fmt"
	"strings"
)

func Test(description string, testFunc func(value any) bool) *Parser {
	return NewParser(description, func(ctx Context) (any, error) {
		if ctx.IsEof() {
			return nil, ParseError{
				Up:            nil,
				Expected:      description,
				Unexpected:    EofString,
				StartLocation: ctx.CurLocation(),
				EndLocation:   ctx.CurLocation(),
			}
		}

		value := ctx.Peek()
		if testFunc(value) {
			return ctx.Pop(), nil
		}
		return nil, ParseError{
			Up:            nil,
			Expected:      description,
			Unexpected:    value,
			StartLocation: ctx.CurLocation(),
			EndLocation:   ctx.CurLocation(),
		}
	})
}

func Exact(expect any) *Parser {
	desc := fmt.Sprintf("exactly '%v'", toOutputAny(expect))
	return Test(desc, func(value any) bool { return value == expect })
}

func Eof() *Parser {
	return NewParser(EofString, func(ctx Context) (any, error) {
		if ctx.IsEof() {
			return nil, nil
		}
		return nil, ParseError{
			Up:            nil,
			Expected:      EofString,
			Unexpected:    ctx.Peek(),
			StartLocation: ctx.CurLocation(),
			EndLocation:   ctx.CurLocation(),
		}
	})
}

func Succeed(value any) *Parser {
	return NewParser("always successful", func(ctx Context) (any, error) {
		return value, nil
	})
}

func Fail(err error) *Parser {
	return NewParser(fmt.Sprintf("always failing with error '%s'", err), func(ctx Context) (any, error) {
		return nil, err
	})
}

func Seq(parsers ...*Parser) *Parser {
	if len(parsers) < 1 {
		panic("Seq parser must have at least 1 parser")
	}

	parserDescriptions := []string{}
	for _, parser := range parsers {
		parserDescriptions = append(parserDescriptions, parser.Description)
	}
	desc := fmt.Sprintf("sequence of %d parsers: ( %s )", len(parsers), strings.Join(parserDescriptions, ", "))

	return NewParser(desc, func(ctx Context) (any, error) {
		results := []any{}
		for _, parser := range parsers {
			startLoc := ctx.CurLocation()
			result, err := parser.Parse(ctx)
			if err != nil {
				return nil, ParseError{
					Up:            err,
					Expected:      desc,
					Unexpected:    GetContextUnexpected(ctx),
					StartLocation: startLoc,
					EndLocation:   ctx.CurLocation(),
				}
			}
			results = append(results, result)
		}
		return results, nil
	})
}

func AnyOf(parsers ...*Parser) *Parser {
	if len(parsers) < 1 {
		panic("Alt parser must have at least 1 parser")
	}

	parserDescriptions := []string{}
	for _, parser := range parsers {
		parserDescriptions = append(parserDescriptions, parser.Description)
	}
	desc := fmt.Sprintf("any of %d parsers: ( %s )", len(parsers), strings.Join(parserDescriptions, ", "))

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
			Unexpected:    GetContextUnexpected(ctx),
			StartLocation: concreteStartLoc,
			EndLocation:   ctx.CurLocation(),
		}
	})
}
