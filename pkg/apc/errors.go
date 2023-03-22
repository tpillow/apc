package apc

import "fmt"

var ErrEOF = &EOFError{}

var ErrParseErr = &ParseError{}

var ErrParseErrConsumed = &ParseErrorConsumed{}

type EOFError struct{}

func (err *EOFError) Error() string {
	return "EOF"
}

func (err *EOFError) Is(target error) bool {
	if _, ok := target.(*EOFError); ok {
		return true
	}
	return false
}

type ParseError struct {
	Err     error
	Message string
	Origin  Origin
}

func ParseErrExpectedButGot(ctx Context, expected interface{}, got interface{}, wrapErr error) *ParseError {
	return &ParseError{
		Err:     wrapErr,
		Message: fmt.Sprintf("expected %v but got %v", expected, got),
		Origin:  ctx.GetCurOrigin(),
	}
}

func ParseErrExpectedButGotNext(ctx Context, expected interface{}, wrapErr error) *ParseError {
	got, _ := ctx.Peek(0, 16) // Note: no error handle here
	if len(got) >= 16 {
		got = fmt.Sprintf("%v ...more...", got)
	}
	return ParseErrExpectedButGot(ctx, expected, got, wrapErr)
}

func (err *ParseError) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("Parse Error at %v: %v", err.Origin, err.Message)
	}
	return fmt.Sprintf("Parse Error at %v: %v\n%v", err.Origin, err.Message, err.Err)
}

func (err *ParseError) Unwrap() error {
	return err.Err
}

func (err *ParseError) Is(target error) bool {
	if _, ok := target.(*ParseError); ok {
		return true
	}
	return false
}

type ParseErrorConsumed struct {
	Err     error
	Message string
	Origin  Origin
}

func ParseErrConsumedExpectedButGot(ctx Context, expected interface{}, got interface{}, wrapErr error) *ParseErrorConsumed {
	return &ParseErrorConsumed{
		Err:     wrapErr,
		Message: fmt.Sprintf("expected %v but got %v", expected, got),
		Origin:  ctx.GetCurOrigin(),
	}
}

func ParseErrConsumedExpectedButGotNext(ctx Context, expected interface{}, wrapErr error) *ParseErrorConsumed {
	got, _ := ctx.Peek(0, 16) // Note: no error handle here
	if len(got) >= 16 {
		got = fmt.Sprintf("%v ...more...", got)
	}
	return ParseErrConsumedExpectedButGot(ctx, expected, got, wrapErr)
}

func (err *ParseErrorConsumed) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("Parse Error (cannot backtrack) at %v: %v", err.Origin, err.Message)
	}
	return fmt.Sprintf("Parse Error (cannot backtrack) at %v: %v\n%v", err.Origin, err.Message, err.Err)
}

func (err *ParseErrorConsumed) Unwrap() error {
	return err.Err
}

func (err *ParseErrorConsumed) Is(target error) bool {
	if _, ok := target.(*ParseErrorConsumed); ok {
		return true
	}
	return false
}
