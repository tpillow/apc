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
		apc.Named("JSON token",
			apc.Any(
				apc.BindToToken(apc.Bind[rune, string, any](apc.ExactStr(string(TokenTypeNull)), nil), TokenTypeNull),
				apc.BindToToken(apc.BoolParser, TokenTypeBool),
				apc.BindToToken(apc.ExactStr(string(TokenTypeColon)), TokenTypeColon),
				apc.BindToToken(apc.ExactStr(string(TokenTypeComma)), TokenTypeComma),
				apc.BindToToken(apc.ExactStr(string(TokenTypeOpenBrace)), TokenTypeOpenBrace),
				apc.BindToToken(apc.ExactStr(string(TokenTypeCloseBrace)), TokenTypeCloseBrace),
				apc.BindToToken(apc.ExactStr(string(TokenTypeOpenBracket)), TokenTypeOpenBracket),
				apc.BindToToken(apc.ExactStr(string(TokenTypeCloseBracket)), TokenTypeCloseBracket),
				apc.BindToToken(apc.DoubleQuotedStringParser, TokenTypeStr),
				apc.BindToToken(apc.FloatParser, TokenTypeNum))),
	)
)

var (
	valueParser    apc.Parser[apc.Token, any]
	valueParserRef = apc.Ref(&valueParser)

	pairParser = apc.Map(
		apc.Seq3(apc.ExactTokenType(TokenTypeStr), apc.ExactTokenType(TokenTypeColon), valueParserRef),
		func(node *apc.Seq3Node[apc.Token, apc.Token, any], _ apc.Origin) PairNode {
			return PairNode{
				Key:   node.Result1.Value.(string),
				Value: node.Result3,
			}
		})

	valueListParser = apc.ZeroOrMoreSeparated(valueParserRef, apc.ExactTokenType(TokenTypeComma))
	pairListParser  = apc.ZeroOrMoreSeparated(pairParser, apc.ExactTokenType(TokenTypeComma))

	objParser = apc.Named("JSON object",
		apc.Map(
			apc.Seq3(apc.ExactTokenType(TokenTypeOpenBrace), pairListParser, apc.ExactTokenType(TokenTypeCloseBrace)),
			func(node *apc.Seq3Node[apc.Token, []PairNode, apc.Token], _ apc.Origin) any {
				return ObjNode{
					Pairs: node.Result2,
				}
			}))

	arrayParser = apc.Named("JSON array",
		apc.Map(
			apc.Seq3(apc.ExactTokenType(TokenTypeOpenBracket), valueListParser, apc.ExactTokenType(TokenTypeCloseBracket)),
			func(node *apc.Seq3Node[apc.Token, []any, apc.Token], _ apc.Origin) ArrayNode {
				return ArrayNode{
					Nodes: node.Result2,
				}
			}))
)

func main() {
	valueParser = apc.Named("JSON value",
		apc.Any(
			apc.MapTokenToValue[apc.Token, any](apc.ExactTokenType(TokenTypeNum)),
			apc.MapTokenToValue[apc.Token, any](apc.ExactTokenType(TokenTypeBool)),
			apc.MapTokenToValue[apc.Token, any](apc.ExactTokenType(TokenTypeNull)),
			apc.MapTokenToValue[apc.Token, any](apc.ExactTokenType(TokenTypeStr)),
			apc.CastToAny(objParser),
			apc.CastToAny(arrayParser)))

	input := ` { "name" : "Tom" , "age" : 55 , "weight":23.35,"hobbies" : [ "sports" , "stuff" , -55, +3.4, [], {} ] } `
	ctx := apc.NewStringContext("<string>", input)
	lexer := apc.NewParseReader[rune](ctx, lexParser)
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
