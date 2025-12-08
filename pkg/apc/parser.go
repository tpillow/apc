package apc

import (
	"fmt"
	"slices"
)

const EofString = "<<EOF>>"

type MapFunc[IT, OT any] func(value IT) OT
type MapSliceFunc[IT, OT any] func(values []IT) OT
type ParserGeneratorFunc[IT, VT, OT any] func(value VT) *Parser[IT, OT]

type ParseFunc[IT, OT any] func(ctx Context[IT]) (OT, error)

// TODO: make interface...
type Parser[IT, OT any] struct {
	Description string
	ParseFunc   ParseFunc[IT, OT]
}

func NewParser[IT, OT any](description string, parseFunc ParseFunc[IT, OT]) *Parser[IT, OT] {
	return &Parser[IT, OT]{
		Description: description,
		ParseFunc:   parseFunc,
	}
}

func NewFutureParser[IT, OT any]() *Parser[IT, OT] {
	return NewParser[IT, OT]("", nil)
}

func (parser *Parser[IT, OT]) Builder() *Builder[IT, OT] {
	return NewBuilder(parser)
}

func (parser *Parser[IT, OT]) Parse(ctx Context[IT]) (OT, error) {
	if parser.ParseFunc == nil {
		panic("cannot use a Parser whose ParseFunc is nil")
	}

	return parser.ParseFunc(ctx)
}

func (parser *Parser[IT, OT]) ParseToEof(ctx Context[IT]) (OT, error) {
	return Skip(parser, Eof[IT]()).Parse(ctx)
}

// TODO: error, what? correct?
func Peek[IT, OT any](parser *Parser[IT, OT]) *Parser[IT, OT] {
	desc := fmt.Sprintf("peeking parser of %s", parser.Description)
	return NewParser(desc, func(ctx Context[IT]) (OT, error) {
		startLoc := ctx.CurLocation()
		result, err := parser.Parse(ctx)
		endLoc := ctx.CurLocation()
		ctx.SetLocation(startLoc)
		return result, ParseError{
			Up:            err,
			Expected:      desc,
			Unexpected:    GetContextUnexpected(ctx),
			StartLocation: startLoc,
			EndLocation:   endLoc,
		}
	})
}

func Map[IT, OT, OT2 any](parser *Parser[IT, OT], transform MapFunc[OT, OT2]) *Parser[IT, OT2] {
	return NewParser(parser.Description, func(ctx Context[IT]) (OT2, error) {
		result, err := parser.Parse(ctx)
		if err != nil {
			return newT[OT2](), err
		}
		return transform(result), nil
	})
}

func Generate[IT, OT, OT2 any](parser *Parser[IT, OT], parserGen ParserGeneratorFunc[IT, OT, OT2]) *Parser[IT, OT2] {
	return NewParser(parser.Description, func(ctx Context[IT]) (OT2, error) {
		result, err := parser.Parse(ctx)
		if err != nil {
			return newT[OT2](), err
		}
		nextResult, err := parserGen(result).Parse(ctx)
		if err != nil {
			return newT[OT2](), err
		}
		return nextResult, err
	})
}

func Index[IT, OT any](parser *Parser[IT, []OT], index int) *Parser[IT, OT] {
	return MapSlice(parser, func(values []OT) OT {
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
func ConcatSlices[IT, OT any](parsers ...*Parser[IT, []OT]) *Parser[IT, []OT] {
	if len(parsers) <= 0 {
		panic("ConcatSlices requires at least 1 parser")
	}

	return MapSlice(Seq(parsers...), func(results [][]OT) []OT {
		return slices.Concat(results...)
	})
}

func Optional[IT, OT any](parser *Parser[IT, OT], defaultValue OT) *Parser[IT, OT] {
	desc := fmt.Sprintf("optional %s", parser.Description)
	return NewParser(desc, func(ctx Context[IT]) (OT, error) {
		startLoc := ctx.CurLocation()
		result, err := parser.Parse(ctx)
		if err != nil {
			ctx.SetLocation(startLoc)
			return defaultValue, nil
		}
		return result, nil
	})
}

func Describe[IT, OT any](parser *Parser[IT, OT], description string) *Parser[IT, OT] {
	// TODO: when parser intf, must change
	parser.Description = description
	return parser
}

func (parser *Parser[IT, OT]) Become(other *Parser[IT, OT]) {
	// TODO: when parser intf, must change
	if parser.ParseFunc != nil {
		panic("cannot call Become on a parser with a non-nil ParseFunc (Become can only be called at most once)")
	}
	parser.Description = other.Description
	parser.ParseFunc = other.ParseFunc
}

func Times[IT, OT any](parser *Parser[IT, OT], min int, max int) *Parser[IT, []OT] {
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

	return NewParser(desc, func(ctx Context[IT]) ([]OT, error) {
		results := []OT{}

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
