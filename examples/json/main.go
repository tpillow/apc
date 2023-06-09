package main

import (
	"fmt"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

var (
	valueParser    apc.Parser[rune, any]
	valueParserRef = apc.Ref(&valueParser)

	pairParser = apc.Named("json key-value pair",
		apc.Map(
			apc.Seq3(apc.DoubleQuotedStringParser, apc.ExactStr(":"), valueParserRef),
			func(node *apc.Seq3Node[string, string, any]) PairNode {
				return PairNode{
					Key:   node.Result1,
					Value: node.Result3,
				}
			}))

	valueListParser = apc.ZeroOrMoreSeparated(valueParserRef, apc.ExactStr(","))
	pairListParser  = apc.ZeroOrMoreSeparated(pairParser, apc.ExactStr(","))

	objParser = apc.Named("JSON object",
		apc.Map(
			apc.Seq3(apc.ExactStr("{"), pairListParser, apc.ExactStr("}")),
			func(node *apc.Seq3Node[string, []PairNode, string]) any {
				return ObjNode{
					Pairs: node.Result2,
				}
			}))

	arrayParser = apc.Named("JSON array",
		apc.Map(
			apc.Seq3(apc.ExactStr("["), valueListParser, apc.ExactStr("]")),
			func(node *apc.Seq3Node[string, []any, string]) ArrayNode {
				return ArrayNode{
					Nodes: node.Result2,
				}
			}))
)

func main() {
	valueParser = apc.Named("JSON value",
		apc.Any(
			apc.CastToAny(apc.FloatParser),
			apc.CastToAny(apc.BoolParser),
			apc.CastToAny(apc.Bind[rune, string, any](apc.ExactStr("null"), nil)),
			apc.CastToAny(apc.DoubleQuotedStringParser),
			apc.CastToAny(objParser),
			apc.CastToAny(arrayParser)))

	input := ` { "name" : "Tom" , "age" : 55 , "weight":23.35,"hobbies" : [ "sports" , "stuff" , -55, +3.4, [], {} ] } `
	ctx := apc.NewStringContext("<string>", input)
	ctx.AddSkipParser(apc.CastToAny(apc.WhitespaceParser))

	fmt.Printf("Input: %v\n", input)
	node, err := apc.Parse[rune](ctx, valueParser, apc.DefaultParseConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Node: %v\n", node)
}

type PairNode struct {
	Key   string
	Value any
}

func (n PairNode) String() string {
	return fmt.Sprintf("Pair<%v: %v>", n.Key, n.Value)
}

type ArrayNode struct {
	Nodes []any
}

func (n ArrayNode) String() string {
	var sb strings.Builder
	sb.WriteString("Array<")
	for i, n := range n.Nodes {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", n))
	}
	sb.WriteString(">")
	return sb.String()
}

type ObjNode struct {
	Pairs []PairNode
}

func (n ObjNode) String() string {
	var sb strings.Builder
	sb.WriteString("Object<")
	for i, n := range n.Pairs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", n))
	}
	sb.WriteString(">")
	return sb.String()
}
