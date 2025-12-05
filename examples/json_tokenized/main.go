package main

import (
	"fmt"

	"github.com/tpillow/apc/pkg/apc"
)

const (
	TokenNull = iota
	TokenInt
	TokenFloat
	TokenStr
	TokenBool
)

type Key string
type Value any
type Dict map[Key]Value

func TokenizeJson(text string) (apc.Token, error) {
	return apc.Token{}, nil
}

func ParseJson(text string) (Value, error) {
	return nil, nil
}

func main() {
	test_json := `{"a": 1, "b": true, "c": [1, 2, 3], "d": {}, "e": []}`
	result, err := ParseJson(test_json)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	fmt.Printf("RESULT: %v\n", result)
}
