package apcgen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

func buildTokenParserForType(buildCtx *buildContext[apc.Token], resultType reflect.Type) apc.Parser[apc.Token, any] {
	return buildParserForTypeCommon(buildCtx, resultType, buildTokenParserFromRootNode)
}

func buildTokenParserFromRootNode(buildCtx *buildContext[apc.Token], subCtx *buildSubcontext[apc.Token], node *rootNode) apc.Parser[apc.Token, any] {
	return buildParserFromRootNodeCommon(buildCtx, subCtx, node, buildTokenParserFromNode)
}

func buildTokenParserFromNode(buildCtx *buildContext[apc.Token], subCtx *buildSubcontext[apc.Token], rawNode Node) apc.Parser[apc.Token, any] {
	if parser := buildParserFromNodeCommon(buildCtx, subCtx, rawNode, buildTokenParserForType, buildTokenParserFromNode); parser != nil {
		return parser
	}
	switch node := rawNode.(type) {
	case *matchStringNode:
		colIdx := strings.Index(node.Value, ":")
		if colIdx == 0 {
			panic(fmt.Sprintf("invalid token 'type' or 'type:value' pair to match token: '%v' "+
				"(if you intend to match a token type of ':' use the explicit matcher: \"token('<type' [, '<value>'])\")", node.Value))
		}
		if colIdx < 0 {
			return apc.Map(
				apc.ExactTokenType(apc.TokenType(node.Value)),
				func(node apc.Token) any {
					if node.Value == nil {
						return node.Type
					}
					return node.Value
				},
			)
		}
		return apc.Map(
			apc.ExactTokenValue(apc.TokenType(node.Value[:colIdx]), node.Value[colIdx+1:]),
			func(node apc.Token) any {
				if node.Value == nil {
					return node.Type
				}
				return node.Value
			},
		)
	case *matchTokenNode:
		var exactTokParser apc.Parser[apc.Token, apc.Token]
		if node.Value.IsNil() {
			exactTokParser = apc.ExactTokenType(apc.TokenType(node.TokenType))
		} else {
			exactTokParser = apc.ExactTokenValue(apc.TokenType(node.TokenType), node.Value.Value())
		}
		return apc.Map(
			exactTokParser,
			func(node apc.Token) any {
				if node.Value == nil {
					return node.Type
				}
				return node.Value
			},
		)
	case *matchRegexNode:
		panic("cannot use 'regex' when using a Token context")
	default:
		panic(fmt.Sprintf("unknown node to process in buildTokenParserFromNode: %T", rawNode))
	}
}
