package main

import (
	"fmt"
	"math"

	"github.com/tpillow/apc/pkg/apc"
)

type Key string
type Value any
type Dict map[Key]Value

type pair struct {
	key   Key
	value Value
}

func ParseJson(text string) (Value, error) {
	keyword_parser := apc.AnyOf(
		apc.ExactStr("true").Builder().MapValue(true).Build(),
		apc.ExactStr("false").Builder().MapValue(false).Build(),
		apc.ExactStr("null").Builder().MapValue(nil).Build())
	simple_value_parser := apc.AnyOf(
		apc.Float64, apc.Int64, apc.DoubleQuotedString, keyword_parser)

	left_bracket_parser := apc.Exact('[')
	right_bracket_parser := apc.Exact(']')
	left_brace_parser := apc.Exact('{')
	right_brace_parser := apc.Exact('}')
	comma_parser := apc.Exact(',')
	colon_parser := apc.Exact(':')

	value_parser := apc.NewFutureParser()
	list_parser := left_bracket_parser.Builder().Then(
		value_parser.Builder().SeparatedBy(comma_parser, 0, math.MaxInt32).Build()).Skip(right_bracket_parser).Describe("array").Build()
	pair_parser := apc.Seq(apc.DoubleQuotedString.Builder().Skip(colon_parser).Build(), value_parser).Builder().Map(func(rawValues any) any {
		values := rawValues.([]any)
		key := values[0].(string)
		return pair{Key(key), Value(values[1])}
	}).Describe("key-value pair").Build()
	dict_parser := left_brace_parser.Builder().Then(
		pair_parser.Builder().SeparatedBy(comma_parser, 0, math.MaxInt32).Build()).Skip(right_brace_parser).Map(func(rawValues any) any {
		values := rawValues.([]any)
		dict := Dict{}
		for _, rawPair := range values {
			pair, ok := rawPair.(pair)
			if !ok {
				panic("must be pair")
			}
			dict[pair.key] = pair.value
		}
		return dict
	}).Describe("object").Build()

	value_parser.Become(
		apc.AnyOf(simple_value_parser, list_parser, dict_parser).Builder().Describe("JSON value").Build())
	parser := value_parser

	ctx := apc.NewStringContext(text)
	ctx.SetPreParser(apc.OptionalWhitespace)
	return parser.ParseToEof(ctx)
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
