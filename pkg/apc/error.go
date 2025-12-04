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
		if err.Up == nil {
			panic("ParseError with no Unexpected should have an Up error")
		}

		baseMsg = fmt.Sprintf("while expecting '%v', got an error at %v (to %v)",
			err.Expected, err.StartLocation, err.EndLocation)
	} else {
		baseMsg = fmt.Sprintf("expected '%v' but got unexpected '%v' at %v (to %v)",
			err.Expected, err.Unexpected, err.StartLocation, err.EndLocation)
	}

	if err.Up == nil {
		return baseMsg
	}

	return fmt.Sprintf("%s\n  stems from:\n%s", baseMsg, err.Up)
}
