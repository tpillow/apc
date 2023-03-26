package apc

import "fmt"

var (
	// Instance of EOFError to compare.
	ErrEOF = &EOFError{}
	// Instance of ParseError to compare.
	ErrParseErr = &ParseError{}
	// Instance of ParseErrorConsumed to compare.
	ErrParseErrConsumed = &ParseErrorConsumed{}
)

func IsMustReturnParseErr(err error) bool {
	if err == nil {
		return false
	}
	if _, ok := err.(*ParseError); ok {
		return false
	}
	return true
}

// EOFError represents that the end of a file or input has been reached.
type EOFError struct{}

// The error string.
func (err *EOFError) Error() string {
	return "EOF"
}

// Returns true if target is also an EOFError.
func (err *EOFError) Is(target error) bool {
	if _, ok := target.(*EOFError); ok {
		return true
	}
	return false
}

// ParseError represents a parser that could not match input and has
// NOT consumed any input.
type ParseError struct {
	// The optional error to wrap.
	Err error
	// The error message.
	Message string
	// The Origin of the error.
	Origin Origin
}

// Returns a ParseError with an error message in the format of "expected but got".
func ParseErrExpectedButGot[CT any](ctx Context[CT], expected interface{}, got interface{}, wrapErr error) *ParseError {
	return &ParseError{
		Err:     wrapErr,
		Message: fmt.Sprintf("expected %v but got %v", expected, got),
		Origin:  ctx.GetCurOrigin(),
	}
}

// Returns a ParseError with an error message in the format of "expected but got" where
// got is the next N input runes (truncated).
func ParseErrExpectedButGotNext[CT any](ctx Context[CT], expected interface{}, wrapErr error) *ParseError {
	got, _ := ctx.Peek(0, 1) // Note: no error handle here
	return ParseErrExpectedButGot(ctx, expected, got, wrapErr)
}

// The error string.
func (err *ParseError) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("Parse Error at %v: %v", err.Origin, err.Message)
	}
	return fmt.Sprintf("Parse Error at %v: %v\n%v", err.Origin, err.Message, err.Err)
}

// Unwraps this error.
func (err *ParseError) Unwrap() error {
	return err.Err
}

// Returns true if target is also a ParseError.
func (err *ParseError) Is(target error) bool {
	if _, ok := target.(*ParseError); ok {
		return true
	}
	return false
}

// ParseErrorConsumed represents a parser that could not match input
// and HAS consumed some input.
type ParseErrorConsumed struct {
	// The optional error to wrap.
	Err error
	// The error message.
	Message string
	// The Origin of the error.
	Origin Origin
}

func ParseErrConsumedExpectedButGot[CT any](ctx Context[CT], expected interface{}, got interface{}, wrapErr error) *ParseErrorConsumed {
	return &ParseErrorConsumed{
		Err:     wrapErr,
		Message: fmt.Sprintf("expected %v but got %v", expected, got),
		Origin:  ctx.GetCurOrigin(),
	}
}

// Returns a ParseError with an error message in the format of "expected but got" where
// got is the next N input runes (truncated).
func ParseErrConsumedExpectedButGotNext[CT any](ctx Context[CT], expected interface{}, wrapErr error) *ParseErrorConsumed {
	got, _ := ctx.Peek(0, 1) // Note: no error handle here
	return ParseErrConsumedExpectedButGot(ctx, expected, got, wrapErr)
}

// The error string.
func (err *ParseErrorConsumed) Error() string {
	if err.Err == nil {
		return fmt.Sprintf("Parse Error (cannot backtrack) at %v: %v", err.Origin, err.Message)
	}
	return fmt.Sprintf("Parse Error (cannot backtrack) at %v: %v\n%v", err.Origin, err.Message, err.Err)
}

// Unwraps this error.
func (err *ParseErrorConsumed) Unwrap() error {
	return err.Err
}

// Returns true if target is also a ParseErrorConsumed.
func (err *ParseErrorConsumed) Is(target error) bool {
	if _, ok := target.(*ParseErrorConsumed); ok {
		return true
	}
	return false
}
