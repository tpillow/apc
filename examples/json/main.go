package main

import (
	"fmt"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

var (
	valueParser    apc.Parser[any]
	valueParserRef = apc.Ref(&valueParser)

	pairParser = apc.Map(
		apc.Seq3("pair", apc.DoubleQuotedStringParser, apc.Exact(":"), valueParserRef),
		func(node apc.Seq3Node[string, string, any]) PairNode {
			return PairNode{
				Key:   node.Result1,
				Value: node.Result3,
			}
		})

	valueListParser = apc.ZeroOrMoreSeparated("value list", valueParserRef, apc.Exact(","))
	pairListParser  = apc.ZeroOrMoreSeparated("pair list", pairParser, apc.Exact(","))

	objParser = apc.Map(
		apc.Seq3("object", apc.Exact("{"), pairListParser, apc.Exact("}")),
		func(node apc.Seq3Node[string, []PairNode, string]) any {
			return ObjNode{
				Pairs: node.Result2,
			}
		})
	arrayParser = apc.Map(
		apc.Seq3("array", apc.Exact("["), valueListParser, apc.Exact("]")),
		func(node apc.Seq3Node[string, []any, string]) ArrayNode {
			return ArrayNode{
				Nodes: node.Result2,
			}
		})
)

func main() {
	valueParser = apc.OneOf("value",
		apc.MapToAny(apc.FloatParser),
		apc.MapToAny(apc.BoolParser),
		apc.MapToAny(apc.Bind[string, any](apc.Exact("null"), nil)),
		apc.MapToAny(apc.DoubleQuotedStringParser),
		apc.MapToAny(objParser),
		apc.MapToAny(arrayParser))

	input := ` { "name" : "Tom" , "age" : 55 , "weight":23.35,"hobbies" : [ "sports" , "stuff" , -55, +3.4 ] } `
	ctx := apc.NewStringContext("<string>", input)
	ctx.AddSkipParser(apc.MapToAny(apc.WhitespaceParser))

	fmt.Printf("Input: %v\n", input)
	node, err := apc.Parse(ctx, valueParser, apc.ParseConfig{MustParseToEOF: true})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Node: %v\n", node)
}

type PairNode struct {
	Key   any
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
