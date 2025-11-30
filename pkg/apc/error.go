package apc

import "fmt"

type ExpErr struct {
	AtIndex    int
	Expected   string
	Unexpected string
}

func (err ExpErr) Error() string {
	return fmt.Sprintf("at index %d: expected '%s' but got '%s'", err.AtIndex, err.Expected, err.Unexpected)
}

type CustomErr struct {
	AtIndex int
	Message string
}

func (err CustomErr) Error() string {
	return fmt.Sprintf("at index %d: %s", err.AtIndex, err.Message)
}
