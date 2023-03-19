package apc

type MapFunc func(node Node) Node

type Mappable interface {
	Map(mapFunc MapFunc) Parser
}

type MapParser struct {
	Parser  Parser
	MapFunc MapFunc
}

func Map(parser Parser, mapFunc MapFunc) *MapParser {
	return &MapParser{
		Parser:  parser,
		MapFunc: mapFunc,
	}
}

func (p *MapParser) Parse(ctx Context) (Node, error) {
	node, err := p.Parser.Parse(ctx)
	if err != nil {
		return nil, err
	}
	return p.MapFunc(node), err
}

func Bind(parser Parser, node Node) Parser {
	return Map(parser, func(_ Node) Node {
		return node
	})
}
