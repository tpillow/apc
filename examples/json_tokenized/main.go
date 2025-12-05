package main

import (
	"fmt"

	"github.com/tpillow/apc/pkg/apc"
)

const (
	TokenNull  apc.TokenType = "null"
	TokenInt   apc.TokenType = "integer"
	TokenFloat apc.TokenType = "float"
	TokenStr   apc.TokenType = "string"
	TokenBool  apc.TokenType = "bool"
)

type Key string
type Value any
type Dict map[Key]Value

func TokenizeJson(text string) (apc.Context, error) {
	ctx := apc.NewStringContext(text)
	rawResults, err := apc.AnyOf().Many().ParseToEof(ctx)
	if err != nil {
		return nil, err
	}

	results := []apc.Token{}
	for _, res := range rawResults {
		results = append(results, res.(Token))
	}
	return result, nil
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
