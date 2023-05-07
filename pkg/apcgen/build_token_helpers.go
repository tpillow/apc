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
		parts := strings.Split(node.Value, ":")
		switch len(parts) {
		case 1:
			return apc.MapTokenToValue[apc.Token, any](apc.ExactTokenType(apc.TokenType(parts[0])))
		case 2:
			return apc.MapTokenToValue[apc.Token, any](apc.ExactTokenValue(apc.TokenType(parts[0]), parts[1]))
		default:
			panic(fmt.Sprintf("unknown token match specifier using a Token context: %v (format is 'TokenType' or 'TokenType:stringValue')", node.Value))
		}
	case *matchRegexNode:
		panic("cannot use 'regex' when using a Token context")
	default:
		panic(fmt.Sprintf("unknown node to process in buildTokenParserFromNode: %T", rawNode))
	}
}