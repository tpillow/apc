package apcgen

import (
	"fmt"
	"reflect"

	"github.com/tpillow/apc/pkg/apc"
)

func buildRuneParserForType(buildCtx *buildContext[rune], resultType reflect.Type) apc.Parser[rune, any] {
	// Return cached parser if available
	if cachedParser := buildCtx.maybeGetCachedParserFromType(resultType); cachedParser != nil {
		return cachedParser
	}
	// If there is a circular reference to a parser not yet generated
	// We must return a placeholder ref parser
	if refParser := buildCtx.maybeMakeRefParserFromType(resultType); refParser != nil {
		return refParser
	}

	// Create subcontext
	subCtx := newBuildSubContextFromType[rune](resultType)
	node, err := parseFull(subCtx.resultTypeElemName, subCtx.grammarText)
	if err != nil {
		panic(fmt.Sprintf("error parsing parser definition for type '%v': %v", subCtx.resultTypeElemName, err))
	}
	parserPtr := new(apc.Parser[rune, any])
	buildCtx.inProgressParserCache[resultType] = parserPtr
	parser := buildRuneParserFromRootNode(buildCtx, subCtx, node)
	buildCtx.parserTypeParserCache[resultType] = parser
	*parserPtr = parser
	delete(buildCtx.inProgressParserCache, resultType)
	return *parserPtr
}

func setRuneCaptureHelper(subCtx *buildSubcontext[rune], result reflect.Value, rawNode any) {
	if arrNode, ok := rawNode.([]any); ok {
		// Recursively process capture results
		for _, elem := range arrNode {
			setRuneCaptureHelper(subCtx, result, elem)
		}
		return
	}

	capNode, ok := rawNode.(captureResult)
	if !ok {
		return
	}

	fieldName := subCtx.fieldNameFromCaptureIdx(capNode.inputIndex)
	field := reflect.Indirect(result).FieldByName(fieldName)
	if !field.IsValid() {
		panic(fmt.Sprintf("field '%v' (%v) not found on type '%v' via reflection", fieldName, field.Kind(), subCtx.resultTypeElemName))
	}

	switch field.Kind() {
	case reflect.Slice:
		fieldElemKind := field.Type().Elem().Kind()
		switch capVal := capNode.value.(type) {
		case []any:
			for _, valElem := range capVal {
				valueSetFieldOrAppendKind(valElem, fieldElemKind, field)
			}
		default:
			valueSetFieldOrAppendKind(capVal, fieldElemKind, field)
		}
	default:
		valueSetFieldOrAppendKind(capNode.value, field.Kind(), field)
	}
}

func buildRuneParserFromRootNode(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], node *rootNode) apc.Parser[rune, any] {
	rootParser := buildRuneParserFromNode(buildCtx, subCtx, node.Child)
	return apc.Named(
		subCtx.resultTypeElemName,
		apc.Map(
			rootParser,
			func(parseNode any) any {
				result := reflect.New(subCtx.resultType.Elem())
				setRuneCaptureHelper(subCtx, result, parseNode)
				return result.Interface()
			},
		),
	)
}

func buildRuneParserFromNode(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], rawNode Node) apc.Parser[rune, any] {
	switch node := rawNode.(type) {
	case *matchStringNode:
		return apc.CastToAny(apc.ExactStr(node.Value))
	case *matchRegexNode:
		return apc.CastToAny(apc.Regex(node.Regex))
	case *inferNode:
		// TODO: if slice, use type of slice
		fieldName := subCtx.fieldNameFromCaptureIdx(node.InputIndex)
		field, ok := subCtx.resultType.Elem().FieldByName(fieldName)
		if !ok {
			panic(fmt.Sprintf("cannot infer parser: field '%v' not found in type '%v'", fieldName, subCtx.resultTypeElemName))
		}
		if field.Type.Kind() == reflect.Slice {
			return buildRuneParserForType(buildCtx, field.Type.Elem())
		}
		return buildRuneParserForType(buildCtx, field.Type)
	case *captureNode:
		return apc.Map(
			buildRuneParserFromNode(buildCtx, subCtx, node.Child),
			func(parseNode any) any {
				return captureResult{
					inputIndex: node.InputIndex,
					value:      parseNode,
				}
			},
		)
	case *seqNode:
		parser := buildSeqRuneParserFromNodes(buildCtx, subCtx, node.Children)
		return apc.CastToAny(parser)
	case *rangeNode:
		childParser := buildRuneParserFromNode(buildCtx, subCtx, node.Child)
		return apc.CastToAny(apc.Range(node.Range.min, node.Range.max, childParser))
	case *orNode:
		return buildAnyRuneParserFromNodes(buildCtx, subCtx, node.Children)
	case *providedParserKeyNode:
		return buildCtx.mustGetProvidedParserByName(node.Name)
	case *maybeNode:
		return apc.CastToAny(apc.Maybe(buildRuneParserFromNode(buildCtx, subCtx, node.Child)))
	default:
		panic(fmt.Sprintf("unknown node to process in buildParserFromNode: %T", rawNode))
	}
}

func buildSeqRuneParserFromNodes(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], nodes []Node) apc.Parser[rune, []any] {
	parsers := make([]apc.Parser[rune, any], len(nodes))
	for i, rawNode := range nodes {
		parsers[i] = buildRuneParserFromNode(buildCtx, subCtx, rawNode)
	}
	return apc.Seq(parsers...)
}

func buildAnyRuneParserFromNodes(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], nodes []Node) apc.Parser[rune, any] {
	parsers := make([]apc.Parser[rune, any], len(nodes))
	for i, rawNode := range nodes {
		parsers[i] = buildRuneParserFromNode(buildCtx, subCtx, rawNode)
	}
	return apc.Any(parsers...)
}
