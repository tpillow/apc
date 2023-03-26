package apc

import (
	"errors"
	"fmt"
)

type TokenType string

type Token struct {
	Type  TokenType
	Value any
}

func (t Token) String() string {
	return fmt.Sprintf("token of type %v ('%v')", t.Type, t.Value)
}

func ExactTokenType(tokenType TokenType) Parser[Token, Token] {
	return func(ctx Context[Token]) (Token, error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return Token{}, err
		}
		vals, err := ctx.Peek(0, 1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return Token{}, err
		}
		val := vals[0]
		if val.Type != tokenType {
			return Token{}, ParseErrExpectedButGot(ctx, fmt.Sprintf("token of type %v", tokenType), val, nil)
		}
		_, err = ctx.Consume(1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return Token{}, err
		}
		return val, nil
	}
}

func ExactTokenValue(tokenType TokenType, value any) Parser[Token, Token] {
	return func(ctx Context[Token]) (Token, error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return Token{}, err
		}
		vals, err := ctx.Peek(0, 1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return Token{}, err
		}
		val := vals[0]
		if val.Type != tokenType || val.Value != value {
			return Token{}, ParseErrExpectedButGot(ctx, fmt.Sprintf("token of type %v ('%v')", tokenType, value), val, nil)
		}
		_, err = ctx.Consume(1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return Token{}, err
		}
		return val, nil
	}
}

func BindTokenType[CT, T any](parser Parser[CT, T], tokenType TokenType) Parser[CT, Token] {
	return Map(parser, func(node T) Token {
		return Token{
			Type:  tokenType,
			Value: node,
		}
	})
}

func MapToTokenValue[CT, T any](parser Parser[CT, Token]) Parser[CT, T] {
	return Map(parser, func(node Token) T {
		return node.Value.(T)
	})
}

func MapToTokenValueAny[CT any](parser Parser[CT, Token]) Parser[CT, any] {
	return MapToTokenValue[CT, any](parser)
}
