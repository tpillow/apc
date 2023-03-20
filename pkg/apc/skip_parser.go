package apc

type SkipParser struct {
	SkipParser Parser
	Parser     Parser
	IsUnskip   bool
}

func Skip(toSkipParser Parser, parser Parser) *SkipParser {
	return &SkipParser{
		SkipParser: toSkipParser,
		Parser:     parser,
		IsUnskip:   false,
	}
}

func Unskip(toUnskipParser Parser, parser Parser) *SkipParser {
	return &SkipParser{
		SkipParser: toUnskipParser,
		Parser:     parser,
		IsUnskip:   true,
	}
}

func (p *SkipParser) Parse(ctx Context) (Node, error) {
	if p.IsUnskip {
		ctx.RemoveSkipParser(p.SkipParser)
		defer ctx.ProcessSkips() // happens after the deferred add below
		defer ctx.AddSkipParser(p.SkipParser)
	} else {
		ctx.AddSkipParser(p.SkipParser)
		defer ctx.RemoveSkipParser(p.SkipParser)
		ctx.ProcessSkips()
	}

	node, err := p.Parser.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (p *SkipParser) Map(mapFunc MapFunc) Parser {
	return Map(p, mapFunc)
}
