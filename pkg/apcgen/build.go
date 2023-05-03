package apcgen

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/tpillow/apc/pkg/apc"
)

func BuildParser[RT any](buildOpts BuildOptions, providedParsers map[string]apc.Parser[rune, any]) apc.Parser[rune, RT] {
	var rtVal RT
	rtType := reflect.TypeOf(rtVal)

	buildCtx := newBuildContext(buildOpts, providedParsers)
	baseParser := buildParserForType(buildCtx, rtType)
	parser := apc.CastTo[rune, any, RT](baseParser)

	if buildCtx.options.SkipWhitespace {
		parser = apc.Skip(
			apc.CastToAny(apc.WhitespaceParser),
			parser,
		)
	}

	return parser
}

func buildParserForType(buildCtx *buildContext[rune], resultType reflect.Type) apc.Parser[rune, any] {
	// Return cached parser if available
	cachedParser := buildCtx.maybeGetCachedParserFromType(resultType)
	if cachedParser != nil {
		return cachedParser
	}

	// Create subcontext
	subCtx := newBuildSubContextFromType[rune](resultType)
	node, err := parseFull(subCtx.resultTypeElemName, subCtx.grammarText)
	if err != nil {
		panic(fmt.Sprintf("error parsing parser definition for type '%v': %v", subCtx.resultTypeElemName, err))
	}
	parser := buildParserFromRootNode(buildCtx, subCtx, node)
	buildCtx.parserTypeParserCache[resultType] = parser
	return parser
}

func setCaptureHelper(subCtx *buildSubcontext[rune], result reflect.Value, rawNode any) {
	if arrNode, ok := rawNode.([]any); ok {
		// Recursively process capture results
		for _, elem := range arrNode {
			setCaptureHelper(subCtx, result, elem)
		}
		return
	}

	capNode, ok := rawNode.(captureResult)
	if !ok {
		return
	}

	fieldName := subCtx.fieldNameFromCaptureIdx(capNode.inputIndex)

	panicHelper := func(reason string) {
		panic(fmt.Sprintf("failed to set field '%v' on type '%v' via reflection when parsing: %v", fieldName, subCtx.resultTypeElemName, reason))
	}

	panicHelperCannotCastTo := func(toType string, val any) {
		panicHelper(fmt.Sprintf("cannot cast to %v: %v", toType, val))
	}

	field := reflect.Indirect(result).FieldByName(fieldName)
	if !field.IsValid() {
		panicHelper("field not found")
	}

	switch field.Kind() {
	case reflect.String:
		if val, ok := capNode.value.(string); ok {
			field.SetString(val)
		} else {
			panicHelperCannotCastTo("string", capNode.value)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val, ok := capNode.value.(string); ok {
			intVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				panicHelperCannotCastTo("int", fmt.Sprintf("%v (%v)", val, err))
			}
			field.SetInt(intVal)
		} else {
			panicHelperCannotCastTo("int (requiring string now)", capNode.value)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val, ok := capNode.value.(string); ok {
			intVal, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				panicHelperCannotCastTo("uint", fmt.Sprintf("%v (%v)", val, err))
			}
			field.SetUint(intVal)
		} else {
			panicHelperCannotCastTo("uint (requiring string now)", capNode.value)
		}
	case reflect.Float32, reflect.Float64:
		if val, ok := capNode.value.(string); ok {
			floatVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				panicHelperCannotCastTo("float", fmt.Sprintf("%v (%v)", val, err))
			}
			field.SetFloat(floatVal)
		} else {
			panicHelperCannotCastTo("float (requiring string now)", capNode.value)
		}
	case reflect.Bool:
		if val, ok := capNode.value.(string); ok {
			boolVal, err := strconv.ParseBool(val)
			if err != nil {
				panicHelperCannotCastTo("bool", fmt.Sprintf("%v (%v)", val, err))
			}
			field.SetBool(boolVal)
		} else {
			panicHelperCannotCastTo("bool (requiring string now)", capNode.value)
		}
	default:
		// TODO: kind for slice, append values...
		if field.CanSet() {
			field.Set(reflect.ValueOf(capNode.value))
		} else {
			panicHelper("field not settable and is not an intrinsic type")
		}
	}
}

func buildParserFromRootNode(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], node *RootNode) apc.Parser[rune, any] {
	rootParser := buildParserFromNode(buildCtx, subCtx, node.Child)
	return apc.Named(
		subCtx.resultTypeElemName,
		apc.Map(
			rootParser,
			func(parseNode any) any {
				result := reflect.New(subCtx.resultType.Elem())
				setCaptureHelper(subCtx, result, parseNode)
				return result.Interface()
			},
		),
	)
}

func buildParserFromNode(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], rawNode Node) apc.Parser[rune, any] {
	switch node := rawNode.(type) {
	case *MatchStringNode:
		return apc.CastToAny(apc.ExactStr(node.Value))
	case *MatchRegexNode:
		return apc.CastToAny(apc.Regex(node.Regex))
	case *InferNode:
		// TODO: if slice, use type of slice
		fieldName := subCtx.fieldNameFromCaptureIdx(node.InputIndex)
		field, ok := subCtx.resultType.Elem().FieldByName(fieldName)
		if !ok {
			panic(fmt.Sprintf("cannot infer parser: field '%v' not found in type '%v'", fieldName, subCtx.resultTypeElemName))
		}
		return buildParserForType(buildCtx, field.Type)
	case *CaptureNode:
		return apc.Map(
			buildParserFromNode(buildCtx, subCtx, node.Child),
			func(parseNode any) any {
				return captureResult{
					inputIndex: node.InputIndex,
					value:      parseNode,
				}
			},
		)
	case *SeqNode:
		parser := buildSeqParserFromNodes(buildCtx, subCtx, node.Children)
		return apc.CastToAny(parser)
	case *RangeNode:
		childParser := buildParserFromNode(buildCtx, subCtx, node.Child)
		return apc.CastToAny(apc.Range(node.Range.Min, node.Range.Max, childParser))
	case *OrNode:
		return buildAnyParserFromNodes(buildCtx, subCtx, node.Children)
	case *ProvidedParserKeyNode:
		return buildCtx.mustGetProvidedParserByName(node.Name)
	default:
		panic(fmt.Sprintf("unknown node to process in buildParserFromNode: %T", rawNode))
	}
}

func buildSeqParserFromNodes(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], nodes []Node) apc.Parser[rune, []any] {
	parsers := make([]apc.Parser[rune, any], len(nodes))
	for i, rawNode := range nodes {
		parsers[i] = buildParserFromNode(buildCtx, subCtx, rawNode)
	}
	return apc.Seq(parsers...)
}

func buildAnyParserFromNodes(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], nodes []Node) apc.Parser[rune, any] {
	parsers := make([]apc.Parser[rune, any], len(nodes))
	for i, rawNode := range nodes {
		parsers[i] = buildParserFromNode(buildCtx, subCtx, rawNode)
	}
	return apc.Any(parsers...)
}
