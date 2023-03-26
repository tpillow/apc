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
		apc.MapToAny(apc.WhitespaceParser),
		apc.OneOf("token",
			apc.BindTokenType(apc.Bind[rune, string, any](apc.ExactStr(string(TokenTypeNull)), nil), TokenTypeNull),
			apc.BindTokenType(apc.BoolParser, TokenTypeBool),
			apc.BindTokenType(apc.ExactStr(string(TokenTypeColon)), TokenTypeColon),
			apc.BindTokenType(apc.ExactStr(string(TokenTypeComma)), TokenTypeComma),
			apc.BindTokenType(apc.ExactStr(string(TokenTypeOpenBrace)), TokenTypeOpenBrace),
			apc.BindTokenType(apc.ExactStr(string(TokenTypeCloseBrace)), TokenTypeCloseBrace),
			apc.BindTokenType(apc.ExactStr(string(TokenTypeOpenBracket)), TokenTypeOpenBracket),
			apc.BindTokenType(apc.ExactStr(string(TokenTypeCloseBracket)), TokenTypeCloseBracket),
			apc.BindTokenType(apc.DoubleQuotedStringParser, TokenTypeStr),
			apc.BindTokenType(apc.FloatParser, TokenTypeNum),
		))
)

var (
	valueParser    apc.Parser[apc.Token, any]
	valueParserRef = apc.Ref(&valueParser)

	pairParser = apc.Map(
		apc.Seq3("pair", apc.ExactTokenType(TokenTypeStr), apc.ExactTokenType(TokenTypeColon), valueParserRef),
		func(node *apc.Seq3Node[apc.Token, apc.Token, any]) PairNode {
			return PairNode{
				Key:   node.Result1.Value.(string),
				Value: node.Result3,
			}
		})

	valueListParser = apc.ZeroOrMoreSeparated("value list", valueParserRef, apc.ExactTokenType(TokenTypeComma))
	pairListParser  = apc.ZeroOrMoreSeparated("pair list", pairParser, apc.ExactTokenType(TokenTypeComma))

	objParser = apc.Map(
		apc.Seq3("object", apc.ExactTokenType(TokenTypeOpenBrace), pairListParser, apc.ExactTokenType(TokenTypeCloseBrace)),
		func(node *apc.Seq3Node[apc.Token, []PairNode, apc.Token]) any {
			return ObjNode{
				Pairs: node.Result2,
			}
		})
	arrayParser = apc.Map(
		apc.Seq3("array", apc.ExactTokenType(TokenTypeOpenBracket), valueListParser, apc.ExactTokenType(TokenTypeCloseBracket)),
		func(node *apc.Seq3Node[apc.Token, []any, apc.Token]) ArrayNode {
			return ArrayNode{
				Nodes: node.Result2,
			}
		})
)

func main() {
	valueParser = apc.OneOf("value",
		apc.MapToTokenValueAny(apc.ExactTokenType(TokenTypeNum)),
		apc.MapToTokenValueAny(apc.ExactTokenType(TokenTypeBool)),
		apc.MapToTokenValueAny(apc.ExactTokenType(TokenTypeNull)),
		apc.MapToTokenValueAny(apc.ExactTokenType(TokenTypeStr)),
		apc.MapToAny(objParser),
		apc.MapToAny(arrayParser))

	input := ` { "name" : "Tom" , "age" : 55 , "weight":23.35,"hobbies" : [ "sports" , "stuff" , -55, +3.4, [], {} ] } `
	ctx := apc.NewStringContext("<string>", input)
	lexer := apc.NewLexer[rune](ctx, lexParser)
	lexerCtx := apc.NewReaderContext[apc.Token](lexer)

	fmt.Printf("Input: %v\n", input)
	node, err := apc.Parse[apc.Token](lexerCtx, valueParser, apc.DefaultParseConfig)
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
