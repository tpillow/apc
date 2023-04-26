package main

import (
	"fmt"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

const (
	TokenTypeStr          apc.TokenType = "string"
	TokenTypeNum          apc.TokenType = "number"
	TokenTypeBool         apc.TokenType = "boolean"
	TokenTypeNull         apc.TokenType = "null"
	TokenTypeColon        apc.TokenType = ":"
	TokenTypeComma        apc.TokenType = ","
	TokenTypeOpenBrace    apc.TokenType = "{"
	TokenTypeCloseBrace   apc.TokenType = "}"
	TokenTypeOpenBracket  apc.TokenType = "["
	TokenTypeCloseBracket apc.TokenType = "]"
)

var (
	lexParser = apc.Skip(
		apc.CastToAny(apc.WhitespaceParser),
		apc.Any("token",
			apc.BindToToken(apc.Bind[rune, string, any](apc.ExactStr(string(TokenTypeNull)), nil), TokenTypeNull),
			apc.CastTokenToAny(apc.BindToToken(apc.BoolParser, TokenTypeBool)),
			apc.CastTokenToAny(apc.BindToToken(apc.ExactStr(string(TokenTypeColon)), TokenTypeColon)),
			apc.CastTokenToAny(apc.BindToToken(apc.ExactStr(string(TokenTypeComma)), TokenTypeComma)),
			apc.CastTokenToAny(apc.BindToToken(apc.ExactStr(string(TokenTypeOpenBrace)), TokenTypeOpenBrace)),
			apc.CastTokenToAny(apc.BindToToken(apc.ExactStr(string(TokenTypeCloseBrace)), TokenTypeCloseBrace)),
			apc.CastTokenToAny(apc.BindToToken(apc.ExactStr(string(TokenTypeOpenBracket)), TokenTypeOpenBracket)),
			apc.CastTokenToAny(apc.BindToToken(apc.ExactStr(string(TokenTypeCloseBracket)), TokenTypeCloseBracket)),
			apc.CastTokenToAny(apc.BindToToken(apc.DoubleQuotedStringParser, TokenTypeStr)),
			apc.CastTokenToAny(apc.BindToToken(apc.FloatParser, TokenTypeNum)),
		))
)

var (
	valueParser    apc.Parser[apc.Token[any], any]
	valueParserRef = apc.Ref(&valueParser)

	pairParser = apc.Map(
		apc.Seq3("pair", apc.ExactTokenType[any](TokenTypeStr), apc.ExactTokenType[any](TokenTypeColon), valueParserRef),
		func(node *apc.Seq3Node[apc.Token[any], apc.Token[any], any], _ apc.Origin) PairNode {
			return PairNode{
				Key:   node.Result1.Value.(string),
				Value: node.Result3,
			}
		})

	valueListParser = apc.ZeroOrMoreSeparated("value list", valueParserRef, apc.ExactTokenType[any](TokenTypeComma))
	pairListParser  = apc.ZeroOrMoreSeparated("pair list", pairParser, apc.ExactTokenType[any](TokenTypeComma))

	objParser = apc.Map(
		apc.Seq3("object", apc.ExactTokenType[any](TokenTypeOpenBrace), pairListParser, apc.ExactTokenType[any](TokenTypeCloseBrace)),
		func(node *apc.Seq3Node[apc.Token[any], []PairNode, apc.Token[any]], _ apc.Origin) any {
			return ObjNode{
				Pairs: node.Result2,
			}
		})
	arrayParser = apc.Map(
		apc.Seq3("array", apc.ExactTokenType[any](TokenTypeOpenBracket), valueListParser, apc.ExactTokenType[any](TokenTypeCloseBracket)),
		func(node *apc.Seq3Node[apc.Token[any], []any, apc.Token[any]], _ apc.Origin) ArrayNode {
			return ArrayNode{
				Nodes: node.Result2,
			}
		})
)

func main() {
	valueParser = apc.Any("value",
		apc.CastToAny(apc.MapTokenToValue(apc.ExactTokenType[any](TokenTypeNum))),
		apc.CastToAny(apc.MapTokenToValue(apc.ExactTokenType[any](TokenTypeBool))),
		apc.CastToAny(apc.MapTokenToValue(apc.ExactTokenType[any](TokenTypeNull))),
		apc.CastToAny(apc.MapTokenToValue(apc.ExactTokenType[any](TokenTypeStr))),
		apc.CastToAny(objParser),
		apc.CastToAny(arrayParser))

	input := ` { "name" : "Tom" , "age" : 55 , "weight":23.35,"hobbies" : [ "sports" , "stuff" , -55, +3.4, [], {} ] } `
	ctx := apc.NewStringContext("<string>", input)
	lexer := apc.NewParseReader[rune](ctx, lexParser)
	lexerCtx := apc.NewReaderContext[apc.Token[any]](lexer)

	fmt.Printf("Input: %v\n", input)
	node, err := apc.Parse[apc.Token[any]](lexerCtx, valueParser, apc.DefaultParseConfig)
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
