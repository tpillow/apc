package apc

import (
	"errors"
	"fmt"
)

// Represents the type of a Token.
type TokenType string

// TokenType representing a nil token.
const NilTokenType TokenType = ""

// Contains a TokenType and some Value.
// If TokenType is NilTokenType, this Token should be assumed to be nil.
type Token struct {
	Type  TokenType
	Value any
}

// Return the string version of a Token.
func (t Token) String() string {
	valStr := ""
	if t.Value != nil {
		valStr = fmt.Sprintf(" ('%v')", t.Value)
	}
	return fmt.Sprintf("token of type %v%v", t.Type, valStr)
}

// Returns a parser that succeeds if the next peeked token from the Context[Token]
// has a Type that is tokenType.
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
		if len(vals) == 0 {
			return Token{}, ParseErrExpectedButGotNext(ctx, fmt.Sprintf("token of type %v", tokenType), nil)
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

// Returns a parser that succeeds if the next peeked token from the Context[Token]
// has a Type that is tokenType and a Value that is value.
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
		if len(vals) == 0 {
			return Token{}, ParseErrExpectedButGotNext(ctx, fmt.Sprintf("token of type %v ('%v')", tokenType, value), nil)
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

// Returns a parser that maps a Parser[CT, T] into a Parser[CT, Token] by returning a new
// Token with Type tokenType and Value being the result of parser.
func BindToToken[CT, T any](parser Parser[CT, T], tokenType TokenType) Parser[CT, Token] {
	return Map(parser, func(node T) Token {
		return Token{
			Type:  tokenType,
			Value: node,
		}
	})
}

// Returns a parser that maps a Parser[CT, Token] into a Parser[CT, T] by returning
// the Value of the parser result.
func MapTokenToValue[CT, T any](parser Parser[CT, Token]) Parser[CT, T] {
	return Map(parser, func(node Token) T {
		if node.Value == nil {
			return zeroVal[T]()
		}
		return node.Value.(T)
	})
}
