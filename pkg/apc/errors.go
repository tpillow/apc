package apc

import "fmt"

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
