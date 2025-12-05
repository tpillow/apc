package apc

import "fmt"

type ParseError struct {
	Up            error
	Expected      any
	Unexpected    any
	StartLocation Location
	EndLocation   Location
}

func toOutputAny(value any) any {
	if strVal, ok := value.(fmt.Stringer); ok {
		return strVal.String()
	}
	if runeVal, ok := value.(rune); ok {
		return string(runeVal)
	}
	return value
}

func (err ParseError) Error() string {
	baseMsg := ""

	if err.StartLocation == err.EndLocation {
		baseMsg += fmt.Sprintf("at %#v:", toOutputAny(err.StartLocation))
	} else {
		baseMsg += fmt.Sprintf("at %#v to %#v:", toOutputAny(err.StartLocation), toOutputAny(err.EndLocation))
	}

	if err.Unexpected == nil {
		baseMsg += fmt.Sprintf("\n  expected: %#v", toOutputAny(err.Expected))
	} else {
		baseMsg += fmt.Sprintf("\n  expected: %#v\n  unexpected: %#v", toOutputAny(err.Expected), toOutputAny(err.Unexpected))
	}

	if err.Up != nil {
		baseMsg += fmt.Sprintf("\n. stems from:\n%s", err.Up.Error())
	}

	return baseMsg
}
