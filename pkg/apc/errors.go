package apc

import "fmt"

// TODO: better way to propagate parse errors vs. IO errors / etc.
// TODO: better error messages for range/one/etc. combinators
type ParseError struct {
	Origin Origin
	Err    error
}

func NewParseError(origin Origin, errFormat string, errFormatArgs ...interface{}) *ParseError {
	return &ParseError{
		Origin: origin,
		Err:    fmt.Errorf(errFormat, errFormatArgs...),
	}
}

func (err *ParseError) Error() string {
	return fmt.Sprintf("Parse Error at %v: %v", err.Origin, err.Err)
}

type EOFError struct{}

func (err *EOFError) Error() string {
	return "EOF"
}
