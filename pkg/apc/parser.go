package apc

import "errors"

type Node interface{}

type Parser interface {
	Parse(ctx Context) (Node, error)
}

type ParseConfig struct {
	MustParseToEOF bool
}

func Parse(ctx Context, parser Parser, config ParseConfig) (Node, error) {
	node, err := parser.Parse(ctx)
	if err != nil {
		return nil, err
	}
	if config.MustParseToEOF {
		r, err := ctx.PeekRune(0)
		if err == nil {
			return nil, NewParseError(ctx.GetOrigin(), "expected EOF but got '%v'", r)
		} else if !errors.Is(err, &EOFError{}) {
			return nil, NewParseError(ctx.GetOrigin(), "expected EOF but got %v", err)
		}
	}
	return node, nil
}
