package apc

import "errors"

type Node interface{}

// TODO: need some type of parser to indicate not to backtrack further (actually doesn't work like that yet)
type Parser interface {
	Parse(ctx Context) (Node, error)
}

type ParseConfig struct {
	MustParseToEOF bool
	SkipParsers    []Parser
}

func Parse(ctx Context, parser Parser, config ParseConfig) (Node, error) {
	for _, p := range config.SkipParsers {
		ctx.AddSkipParser(p)
	}

	ctx.ProcessSkips()
	node, err := parser.Parse(ctx)
	if err != nil {
		return nil, err
	}
	ctx.ProcessSkips()

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
