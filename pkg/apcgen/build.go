package apcgen

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/tpillow/apc/pkg/apc"
)

func BuildParser[RT any](buildOpts BuildOptions, providedParsers map[string]apc.Parser[rune, any]) apc.Parser[rune, *RT] {
	rtType := reflect.TypeOf(new(RT))

	buildCtx := newBuildContext(buildOpts, providedParsers)
	baseParser := buildParserForType(buildCtx, rtType)
	parser := apc.CastTo[rune, any, *RT](baseParser)

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
	parser := buildParserFromRootNode(buildCtx, subCtx, node)
	buildCtx.parserTypeParserCache[resultType] = parser
	*parserPtr = parser
	delete(buildCtx.inProgressParserCache, resultType)
	return *parserPtr
}

func valueSetFieldOrAppendKind(rawVal any, valKind reflect.Kind, field reflect.Value) {
	maybeAppendValToSliceTrueIfNot := func(val any) bool {
		if field.Kind() != reflect.Slice {
			return true
		}
		field.Set(reflect.Append(field, reflect.ValueOf(val)))
		return false
	}

	panicUnsettable := func(val any, exp string) {
		panic(fmt.Sprintf("cannot set field to value '%v': cannot convert %v", val, exp))
	}

	switch valKind {
	case reflect.String:
		switch val := rawVal.(type) {
		case string:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetString(val)
			}
		default:
			panicUnsettable(rawVal, "to string")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch val := rawVal.(type) {
		case int, int8, int16, int32, int64:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetInt(val.(int64))
			}
		case string:
			if cVal, err := strconv.ParseInt(val, 10, 64); err == nil {
				if maybeAppendValToSliceTrueIfNot(cVal) {
					field.SetInt(cVal)
				}
			} else {
				panicUnsettable(rawVal, "to int from string")
			}
		default:
			panicUnsettable(rawVal, "to int")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch val := rawVal.(type) {
		case uint, uint8, uint16, uint32, uint64:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetUint(val.(uint64))
			}
		case string:
			if cVal, err := strconv.ParseUint(val, 10, 64); err == nil {
				if maybeAppendValToSliceTrueIfNot(cVal) {
					field.SetUint(cVal)
				}
			} else {
				panicUnsettable(rawVal, "to uint from string")
			}
		default:
			panicUnsettable(rawVal, "to uint")
		}
	case reflect.Float32, reflect.Float64:
		switch val := rawVal.(type) {
		case float32, float64:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetFloat(val.(float64))
			}
		case string:
			if cVal, err := strconv.ParseFloat(val, 64); err == nil {
				if maybeAppendValToSliceTrueIfNot(cVal) {
					field.SetFloat(cVal)
				}
			} else {
				panicUnsettable(rawVal, "to float from string")
			}
		default:
			panicUnsettable(rawVal, "to float")
		}
	case reflect.Bool:
		switch val := rawVal.(type) {
		case bool:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetBool(val)
			}
		case string:
			if cVal, err := strconv.ParseBool(val); err == nil {
				if maybeAppendValToSliceTrueIfNot(cVal) {
					field.SetBool(cVal)
				}
			} else {
				panicUnsettable(rawVal, "to bool from string")
			}
		case apc.MaybeValue[any]:
			field.SetBool(!val.IsNil())
		case apc.MaybeValue[string]:
			field.SetBool(!val.IsNil())
		case apc.MaybeValue[apc.Token]:
			field.SetBool(!val.IsNil())
		default:
			panicUnsettable(rawVal, "to bool")
		}
	case reflect.Pointer:
		switch val := rawVal.(type) {
		case apc.MaybeValue[any]:
			if !val.IsNil() {
				if maybeAppendValToSliceTrueIfNot(val.Value()) {
					field.Set(reflect.ValueOf(val.Value()))
				}
			}
		default:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.Set(reflect.ValueOf(val))
			}
		}
	default:
		panicUnsettable(rawVal, fmt.Sprintf("unsupported value kind %v", valKind))
	}
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

func buildParserFromRootNode(buildCtx *buildContext[rune], subCtx *buildSubcontext[rune], node *rootNode) apc.Parser[rune, any] {
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
			return buildParserForType(buildCtx, field.Type.Elem())
		}
		return buildParserForType(buildCtx, field.Type)
	case *captureNode:
		return apc.Map(
			buildParserFromNode(buildCtx, subCtx, node.Child),
			func(parseNode any) any {
				return captureResult{
					inputIndex: node.InputIndex,
					value:      parseNode,
				}
			},
		)
	case *seqNode:
		parser := buildSeqParserFromNodes(buildCtx, subCtx, node.Children)
		return apc.CastToAny(parser)
	case *rangeNode:
		childParser := buildParserFromNode(buildCtx, subCtx, node.Child)
		return apc.CastToAny(apc.Range(node.Range.min, node.Range.max, childParser))
	case *orNode:
		return buildAnyParserFromNodes(buildCtx, subCtx, node.Children)
	case *providedParserKeyNode:
		return buildCtx.mustGetProvidedParserByName(node.Name)
	case *maybeNode:
		return apc.CastToAny(apc.Maybe(buildParserFromNode(buildCtx, subCtx, node.Child)))
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
