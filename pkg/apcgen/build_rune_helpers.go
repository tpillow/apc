package apcgen

import (
	"fmt"
	"reflect"

	"github.com/tpillow/apc/pkg/apc"
)

func buildRuneParserForType(buildCtx *buildContext[rune], resultType reflect.Type) apc.Parser[rune, any] {
	return buildParserForTypeCommon(buildCtx, resultType, buildRuneParserFromRootNode)
}

func buildRuneParserFromRootNode(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], node *rootNode) apc.Parser[rune, any] {
	return buildParserFromRootNodeCommon(buildCtx, subCtx, node, buildRuneParserFromNode)
}

func buildRuneParserFromNode(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], rawNode Node) apc.Parser[rune, any] {
	if parser := buildParserFromNodeCommon(buildCtx, subCtx, rawNode, buildRuneParserForType, buildRuneParserFromNode); parser != nil {
		return parser
	}
	switch node := rawNode.(type) {
	case *matchStringNode:
		return apc.CastToAny(apc.ExactStr(node.Value))
	case *matchRegexNode:
		return apc.CastToAny(apc.Regex(node.Regex))
	case *matchTokenNode:
		panic("cannot use 'token' when using a Token context")
	default:
		panic(fmt.Sprintf("unknown node to process in buildRuneParserFromNode: %T", rawNode))
	}
}
