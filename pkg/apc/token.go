package apc

type TokenType int

type Token struct {
	Type          TokenType
	Value         any
	StartLocation Location
	EndLocation   Location
}

func (parser *Parser) MapToken(tokenType TokenType) *Parser {
	return NewParser(parser.Description, func(ctx Context) (any, error) {
		startLoc := ctx.CurLocation()
		result, err := parser.Parse(ctx)
		if err != nil {
			return nil, err
		}
		return Token{
			Type:          tokenType,
			Value:         result,
			StartLocation: startLoc,
			EndLocation:   ctx.CurLocation(),
		}, nil
	})
}
