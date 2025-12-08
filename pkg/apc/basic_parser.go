package apc

import (
	"fmt"
	"strings"
)

func Test[IT any](description string, testFunc func(value IT) bool) *Parser[IT, IT] {
	return NewParser(description, func(ctx Context[IT]) (IT, error) {
		if ctx.IsEof() {
			return newT[IT](), ParseError{
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
		return newT[IT](), ParseError{
			Up:            nil,
			Expected:      description,
			Unexpected:    value,
			StartLocation: ctx.CurLocation(),
			EndLocation:   ctx.CurLocation(),
		}
	})
}

func Exact[IT comparable](expect IT) *Parser[IT, IT] {
	desc := fmt.Sprintf("exactly '%v'", toOutputAny(expect))
	return Test(desc, func(value IT) bool { return value == expect })
}

func Eof[IT any]() *Parser[IT, IT] {
	return NewParser(EofString, func(ctx Context[IT]) (IT, error) {
		if ctx.IsEof() {
			return newT[IT](), nil
		}
		return newT[IT](), ParseError{
			Up:            nil,
			Expected:      EofString,
			Unexpected:    ctx.Peek(),
			StartLocation: ctx.CurLocation(),
			EndLocation:   ctx.CurLocation(),
		}
	})
}

func Succeed[IT, OT any](value OT) *Parser[IT, OT] {
	return NewParser("always successful", func(ctx Context[IT]) (OT, error) {
		return value, nil
	})
}

func Fail[IT any](err error) *Parser[IT, IT] {
	return NewParser(fmt.Sprintf("always failing with error '%s'", err), func(ctx Context[IT]) (IT, error) {
		return newT[IT](), err
	})
}

func Seq[IT, OT any](parsers ...*Parser[IT, OT]) *Parser[IT, []OT] {
	if len(parsers) < 1 {
		panic("Seq parser must have at least 1 parser")
	}

	parserDescriptions := []string{}
	for _, parser := range parsers {
		parserDescriptions = append(parserDescriptions, parser.Description)
	}
	desc := fmt.Sprintf("sequence of %d parsers: ( %s )", len(parsers), strings.Join(parserDescriptions, ", "))

	return NewParser(desc, func(ctx Context[IT]) ([]OT, error) {
		results := []OT{}
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

type Result2[T1, T2 any] struct {
	Value1 T1
	Value2 T2
}

func Seq2[IT, OT, OT2 any](parser1 *Parser[IT, OT], parser2 *Parser[IT, OT2]) *Parser[IT, Result2[OT, OT2]] {
	parserDescriptions := []string{
		parser1.Description, parser2.Description,
	}
	desc := fmt.Sprintf("sequence of 2 parsers: ( %s )", strings.Join(parserDescriptions, ", "))

	return NewParser(desc, func(ctx Context[IT]) (Result2[OT, OT2], error) {
		startLoc := ctx.CurLocation()
		result1, err := parser1.Parse(ctx)
		if err != nil {
			return Result2[OT, OT2]{}, ParseError{
				Up:            err,
				Expected:      desc,
				Unexpected:    GetContextUnexpected(ctx),
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}

		startLoc = ctx.CurLocation()
		result2, err := parser2.Parse(ctx)
		if err != nil {
			return Result2[OT, OT2]{}, ParseError{
				Up:            err,
				Expected:      desc,
				Unexpected:    GetContextUnexpected(ctx),
				StartLocation: startLoc,
				EndLocation:   ctx.CurLocation(),
			}
		}

		return Result2[OT, OT2]{
			Value1: result1,
			Value2: result2,
		}, nil
	})
}

func AnyOf[IT, OT any](parsers ...*Parser[IT, OT]) *Parser[IT, OT] {
	if len(parsers) < 1 {
		panic("Alt parser must have at least 1 parser")
	}

	parserDescriptions := []string{}
	for _, parser := range parsers {
		parserDescriptions = append(parserDescriptions, parser.Description)
	}
	desc := fmt.Sprintf("any of %d parsers: ( %s )", len(parsers), strings.Join(parserDescriptions, ", "))

	return NewParser(desc, func(ctx Context[IT]) (OT, error) {
		var err error
		concreteStartLoc := ctx.CurLocation()

		for i, parser := range parsers {
			startLoc := ctx.CurLocation()
			var result OT
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

		return newT[OT](), ParseError{
			Up:            err,
			Expected:      desc,
			Unexpected:    GetContextUnexpected(ctx),
			StartLocation: concreteStartLoc,
			EndLocation:   ctx.CurLocation(),
		}
	})
}
