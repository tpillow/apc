package apc

import "fmt"

type TokenType string

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

func (parser *Parser) ExactTokenType(tokenType TokenType) *Parser {
	desc := fmt.Sprintf("token of type %s", tokenType)
	return Test(desc, func(rawToken any) bool {
		token, ok := rawToken.(Token)
		if !ok {
			panic("cannot use ExactTokenType on a non-Token stream")
		}
		return token.Type == tokenType
	})
}

func (parser *Parser) ExactTokenValue(tokenValue any) *Parser {
	desc := fmt.Sprintf("token with value %#v", toOutputAny(tokenValue))
	return Test(desc, func(rawToken any) bool {
		token, ok := rawToken.(Token)
		if !ok {
			panic("cannot use ExactTokenValue on a non-Token stream")
		}
		return token.Value == tokenValue
	})
}

func (parser *Parser) ExactToken(tokenType TokenType, tokenValue any) *Parser {
	desc := fmt.Sprintf("token of type %s with value %#v", tokenType, toOutputAny(tokenValue))
	return Test(desc, func(rawToken any) bool {
		token, ok := rawToken.(Token)
		if !ok {
			panic("cannot use ExactToken on a non-Token stream")
		}
		return token.Type == tokenType && token.Value == tokenValue
	})
}
