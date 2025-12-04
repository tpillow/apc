package apc

import "fmt"

type ParseError struct {
	Up            error
	Expected      any
	Unexpected    any
	StartLocation Location
	EndLocation   Location
}

func (err ParseError) Error() string {
	baseMsg := ""
	if err.Unexpected == nil {
		baseMsg = fmt.Sprintf("at %v to %v:\n  expected: %v",
			err.StartLocation, err.EndLocation, err.Expected)
	} else {
		baseMsg = fmt.Sprintf("at %v to %v:\n  expected: %v\n  unexpected: %v",
			err.StartLocation, err.EndLocation, err.Expected, err.Unexpected)
	}

	if err.Up == nil {
		return baseMsg
	}

	return fmt.Sprintf("%s\n  error stems from:\n%s", baseMsg, err.Up)
}
