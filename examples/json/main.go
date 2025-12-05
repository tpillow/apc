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
		apc.ExactStr("true").Bind(true),
		apc.ExactStr("false").Bind(false),
		apc.ExactStr("null").Bind(nil))
	simple_value_parser := apc.AnyOf(
		apc.Float, apc.Int, apc.DoubleQuotedString, keyword_parser)

	left_bracket_parser := apc.Exact('[').Skip(apc.OptionalWhitespace)
	right_bracket_parser := apc.Exact(']').Skip(apc.OptionalWhitespace)
	left_brace_parser := apc.Exact('{').Skip(apc.OptionalWhitespace)
	right_brace_parser := apc.Exact('}').Skip(apc.OptionalWhitespace)
	comma_parser := apc.Exact(',').Skip(apc.OptionalWhitespace)
	colon_parser := apc.Exact(':').Skip(apc.OptionalWhitespace)

	value_parser := apc.NewFutureParser()
	list_parser := left_bracket_parser.Then(
		value_parser.SeparatedBy(comma_parser, 0, math.MaxInt32)).Skip(right_bracket_parser)
	pair_parser := apc.Seq(apc.DoubleQuotedString.Skip(colon_parser), value_parser).Map(func(rawValues any) any {
		values := rawValues.([]any)
		key := values[0].(Key)
		val := values[1].(Value)
		return pair{key, val}
	})
	dict_parser := left_brace_parser.Then(
		pair_parser.SeparatedBy(comma_parser, 0, math.MaxInt32)).Skip(right_brace_parser).Map(func(rawPairs any) any {
		pairs := rawPairs.([]pair)
		dict := Dict{}
		for _, pair := range pairs {
			dict[pair.key] = pair.value
		}
		return dict
	})

	value_parser.Become(
		apc.AnyOf(simple_value_parser, list_parser, dict_parser).Skip(apc.OptionalWhitespace).Describe("JSON value"))
	parser := apc.OptionalWhitespace.Then(value_parser).Skip(apc.OptionalWhitespace)

	ctx := apc.NewStringContext(text)
	return parser.ParseToEof(ctx)
}

func main() {
	test_json := `"hi"`
	result, err := ParseJson(test_json)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RESULT: %v\n", result)
}
