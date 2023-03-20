package apc

type OneParser struct {
	Parsers []Parser
}

func One(parsers ...Parser) *OneParser {
	return &OneParser{
		Parsers: parsers,
	}
}

func (p *OneParser) Parse(ctx Context) (Node, error) {
	for _, parser := range p.Parsers {
		ctx.ProcessSkips()
		node, err := parser.Parse(ctx)
		if err != nil {
			continue
		}
		return node, nil
	}
	return nil, NewParseError(ctx.GetOrigin(), "expected any but got none")
}

func (p *OneParser) Map(mapFunc MapFunc) Parser {
	return Map(p, mapFunc)
}
