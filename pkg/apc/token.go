package apc

import (
	"errors"
	"fmt"
)

type TokenType string

type Token[T any] struct {
	Type  TokenType
	Value T
}

func (t Token[T]) String() string {
	return fmt.Sprintf("token of type %v ('%v')", t.Type, t.Value)
}

func (t Token[T]) IsNil() bool {
	return t.Type == ""
}

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

func BindToToken[CT, T any](parser Parser[CT, T], tokenType TokenType) Parser[CT, Token[T]] {
	return Map(parser, func(node T, _ Origin) Token[T] {
		return Token[T]{
			Type:  tokenType,
			Value: node,
		}
	})
}

func MapTokenToValue[CT, T any](parser Parser[CT, Token[T]]) Parser[CT, T] {
	return Map(parser, func(node Token[T], _ Origin) T {
		return node.Value
	})
}

func CastTokenTo[CT, T, U any](parser Parser[CT, Token[T]]) Parser[CT, Token[U]] {
	return Map(parser, func(node Token[T], _ Origin) Token[U] {
		return Token[U]{
			Type:  node.Type,
			Value: any(node.Value).(U),
		}
	})
}

func CastTokenToAny[CT, T any](parser Parser[CT, Token[T]]) Parser[CT, Token[any]] {
	return Map(parser, func(node Token[T], _ Origin) Token[any] {
		return Token[any]{
			Type:  node.Type,
			Value: any(node.Value),
		}
	})
}
