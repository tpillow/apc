package apc

type ExactParser struct {
	Value string
}

func Exact(value string) *ExactParser {
	return &ExactParser{
		Value: value,
	}
}

func (p *ExactParser) Parse(ctx Context) (Node, error) {
	peek, err := PeekNRunes(ctx, 0, len(p.Value))
	if err != nil {
		return nil, NewParseError(ctx.GetOrigin(), "expected '%v' but got '%v'", p.Value, peek)
	}
	if peek != p.Value {
		return nil, NewParseError(ctx.GetOrigin(), "expected '%v' but got '%v'", p.Value, peek)
	}
	val, err := ConsumeNRunes(ctx, len(p.Value))
	if err != nil {
		return nil, NewParseError(ctx.GetOrigin(), "expected '%v' but got '%v'", p.Value, peek)
	}
	return val, nil
}

func (p *ExactParser) Map(mapFunc MapFunc) Parser {
	return Map(p, mapFunc)
}
