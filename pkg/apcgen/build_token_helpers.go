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
		colIdx := strings.LastIndex(node.Value, ":")
		if colIdx <= 0 {
			// case of 0 == ':' maps to a token type of ':'
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
	case *matchRegexNode:
		panic("cannot use 'regex' when using a Token context")
	default:
		panic(fmt.Sprintf("unknown node to process in buildTokenParserFromNode: %T", rawNode))
	}
}
