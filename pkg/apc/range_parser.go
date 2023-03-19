package apc

type RangeParser struct {
	Parser Parser
	Min    int
	Max    int
}

func Range(parser Parser, min int, max int) *RangeParser {
	return &RangeParser{
		Parser: parser,
		Min:    min,
		Max:    max,
	}
}

func ZeroOrOne(parser Parser) Parser {
	return Range(parser, 0, 1).Map(func(node Node) Node {
		nodes := node.([]Node)
		if nodes == nil || len(nodes) <= 0 {
			return nil
		}
		if len(nodes) != 1 {
			panic("unreachable")
		}
		return nodes[0]
	})
}

func ZeroOrMore(parser Parser) Parser {
	return Range(parser, 0, -1)
}

func OneOrMore(parser Parser) Parser {
	return Range(parser, 1, -1)
}

func (p *RangeParser) Parse(ctx Context) (Node, error) {
	nodes := make([]Node, 0)
	for node, err := p.Parser.Parse(ctx); err == nil; {
		nodes = append(nodes, node)
		node, err = p.Parser.Parse(ctx)
	}
	if p.Min >= 0 && len(nodes) < p.Min {
		// TODO: better error
		return nil, NewParseError(ctx.GetOrigin(), "expected at least %v but got %v", p.Min, len(nodes))
	}
	if p.Max >= 0 && len(nodes) > p.Max {
		// TODO: better error
		return nil, NewParseError(ctx.GetOrigin(), "expected at most %v but got %v", p.Max, len(nodes))
	}
	return nodes, nil
}

func (p *RangeParser) Map(mapFunc MapFunc) Parser {
	return Map(p, mapFunc)
}
