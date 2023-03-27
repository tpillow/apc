package apc

import (
	"errors"
	"fmt"
)

// Represents the type of a Token.
type TokenType string

// Contains a TokenType and some Value.
type Token[T any] struct {
	Type  TokenType
	Value T
}

// Return the string version of a Token[T].
func (t Token[T]) String() string {
	return fmt.Sprintf("token of type %v ('%v')", t.Type, t.Value)
}

// Returns true if the Token[T] is "nil" (the default token).
// This is true if the token's Type is "".
func (t Token[T]) IsNil() bool {
	return t.Type == ""
}

// Returns a parser that succeeds if the next peeked token from the Context[Token[T]]
// has a Type that is tokenType.
func ExactTokenType[T any](tokenType TokenType) Parser[Token[T], Token[T]] {
	return func(ctx Context[Token[T]]) (Token[T], error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return Token[T]{}, err
		}
		vals, err := ctx.Peek(0, 1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return Token[T]{}, err
		}
		val := vals[0]
		if val.Type != tokenType {
			return Token[T]{}, ParseErrExpectedButGot(ctx, fmt.Sprintf("token of type %v", tokenType), val, nil)
		}
		_, err = ctx.Consume(1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return Token[T]{}, err
		}
		return val, nil
	}
}

// Returns a parser that succeeds if the next peeked token from the Context[Token[T]]
// has a Type that is tokenType and a Value that is value.
func ExactTokenValue[T any](tokenType TokenType, value T) Parser[Token[T], Token[T]] {
	return func(ctx Context[Token[T]]) (Token[T], error) {
		err := ctx.RunSkipParsers()
		if err != nil {
			return Token[T]{}, err
		}
		vals, err := ctx.Peek(0, 1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return Token[T]{}, err
		}
		val := vals[0]
		if val.Type != tokenType || any(val.Value) != any(value) {
			return Token[T]{}, ParseErrExpectedButGot(ctx, fmt.Sprintf("token of type %v ('%v')", tokenType, value), val, nil)
		}
		_, err = ctx.Consume(1)
		if err != nil && !errors.Is(err, ErrEOF) {
			return Token[T]{}, err
		}
		return val, nil
	}
}

// Returns a parser that maps a Parser[CT, T] into a Parser[CT, Token[T]] by returning a new
// Token[T] with Type tokenType and Value being the result of parser.
func BindToToken[CT, T any](parser Parser[CT, T], tokenType TokenType) Parser[CT, Token[T]] {
	return Map(parser, func(node T, _ Origin) Token[T] {
		return Token[T]{
			Type:  tokenType,
			Value: node,
		}
	})
}

// Returns a parser that maps a Parser[CT, Token[T]] into a Parser[CT, T] by returning
// the Value of the parser result.
func MapTokenToValue[CT, T any](parser Parser[CT, Token[T]]) Parser[CT, T] {
	return Map(parser, func(node Token[T], _ Origin) T {
		return node.Value
	})
}

// Returns a parser that maps a Parser[CT, Token[T]] into a Parser[CT, Token[U]] by returning
// a new Token[U] with Type being the Type of the result of the parser, and Value being
// the Value of the result of the parser casted to type U.
func CastTokenTo[CT, T, U any](parser Parser[CT, Token[T]]) Parser[CT, Token[U]] {
	return Map(parser, func(node Token[T], _ Origin) Token[U] {
		return Token[U]{
			Type:  node.Type,
			Value: any(node.Value).(U),
		}
	})
}

// Equivalent to CastTokenTo[CT, T, any](parser).
func CastTokenToAny[CT, T any](parser Parser[CT, Token[T]]) Parser[CT, Token[any]] {
	return Map(parser, func(node Token[T], _ Origin) Token[any] {
		return Token[any]{
			Type:  node.Type,
			Value: any(node.Value),
		}
	})
}
