package apcgen

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/kr/pretty"
	"github.com/tpillow/apc/pkg/apc"
)

func wrapWithSkipParsers[CT, RT any](parser apc.Parser[CT, RT], skipParsers []apc.Parser[CT, any]) apc.Parser[CT, RT] {
	for _, skipParser := range skipParsers {
		parser = apc.Skip(skipParser, parser)
	}
	return parser
}

func buildParserForTypeCommon[CT any](buildCtx *buildContext[CT], resultType reflect.Type,
	buildParserFromRootNodeFunc func(*buildContext[CT], *buildSubcontext[CT], *rootNode) apc.Parser[CT, any]) apc.Parser[CT, any] {
	// Return cached parser if available
	if cachedParser := buildCtx.parserCache.maybeGetCachedParserFromType(resultType); cachedParser != nil {
		return cachedParser
	}

	// Create subcontext
	subCtx := newBuildSubContextFromType[CT](resultType)

	// Parse the grammar of the result type struct
	node, err := parseFull(subCtx.resultStructType.Name(), subCtx.grammarText, false)
	if err != nil {
		locStr := ""
		if parseErr, ok := err.(*apc.ParseError); ok {
			locStr = "\n" + strings.Repeat(" ", parseErr.Origin.ColNum-1) + "^"
		}
		panic(fmt.Sprintf("error parsing parser definition for type '%v': %v\n%v%v",
			subCtx.resultStructType.Name(), err, subCtx.grammarText, locStr))
	}
	// Debug print the built parser
	maybeLog(DebugPrintBuiltNodes, "Built parser of type %v: %v", subCtx.resultStructType.Name(), pretty.Sprint(node))

	// Log that this result type parser is being built
	parserPtr := new(apc.Parser[CT, any])
	buildCtx.parserCache[resultType] = parserPtr
	// Actually build the parser and set it
	*parserPtr = buildParserFromRootNodeFunc(buildCtx, subCtx, node)
	return *parserPtr
}

func buildParserFromRootNodeCommon[CT any](buildCtx *buildContext[CT], subCtx *buildSubcontext[CT], node *rootNode,
	buildParserFromNodeFunc func(*buildContext[CT], *buildSubcontext[CT], Node) apc.Parser[CT, any]) apc.Parser[CT, any] {
	rootParser := buildParserFromNodeFunc(buildCtx, subCtx, node.Child)
	return apc.Named(
		subCtx.resultStructType.Name(),
		apc.MapDetailed(
			rootParser,
			func(parseNode any, orgRange apc.OriginRange) (any, error) {
				// Set fields of the resulting populated struct
				resultPtrVal := reflect.New(subCtx.resultStructType)
				setCaptureHelper(subCtx, resultPtrVal, parseNode)

				// If the struct has an "OriginRange" field of type apc.OriginRange, set the origin range that matched
				if orgRangeField, ok := subCtx.resultStructType.FieldByName("OriginRange"); ok && orgRangeField.Type == reflectTypeOf[apc.OriginRange]() {
					resultPtrVal.Elem().FieldByName("OriginRange").Set(reflect.ValueOf(orgRange))
				}

				// Return the actual value
				switch subCtx.resultType.Kind() {
				case reflect.Pointer:
					return resultPtrVal.Interface(), nil
				default:
					return reflect.Indirect(resultPtrVal).Interface(), nil
				}
			},
		),
	)
}

func setCaptureHelper[CT any](subCtx *buildSubcontext[CT], resultPtrVal reflect.Value, rawVal any) {
	if arrNode, ok := rawVal.([]any); ok {
		// Recursively process capture results
		for _, elem := range arrNode {
			setCaptureHelper(subCtx, resultPtrVal, elem)
		}
		return
	}

	// TODO: validate this...
	if maybeNode, ok := rawVal.(apc.MaybeValue[[]any]); ok {
		if !maybeNode.IsNil() {
			for _, elem := range maybeNode.Value() {
				setCaptureHelper(subCtx, resultPtrVal, elem)
			}
		}
	} else if maybeNode, ok := rawVal.(apc.MaybeValue[any]); ok {
		if !maybeNode.IsNil() {
			setCaptureHelper(subCtx, resultPtrVal, maybeNode.Value())
		}
	}

	capNode, ok := rawVal.(captureResult)
	if !ok {
		return
	}

	fieldName := subCtx.fieldNameFromCaptureIdx(capNode.inputIndex)
	field := reflect.Indirect(resultPtrVal).FieldByName(fieldName)
	if !field.IsValid() {
		panic(fmt.Sprintf("field '%v' (%v) not found on type '%v' via reflection", fieldName, field.Kind(), subCtx.resultStructType.Name()))
	}

	switch field.Kind() {
	case reflect.Slice:
		fieldElemKind := field.Type().Elem().Kind()
		switch capVal := capNode.value.(type) {
		case []any:
			for _, valElem := range capVal {
				valueSetFieldOrAppendKind(valElem, fieldElemKind, fieldName, field)
			}
		default:
			valueSetFieldOrAppendKind(capVal, fieldElemKind, fieldName, field)
		}
	default:
		valueSetFieldOrAppendKind(capNode.value, field.Kind(), fieldName, field)
	}
}

func buildSeqParserFromNodesCommon[CT any](buildCtx *buildContext[CT], subCtx *buildSubcontext[CT], nodes []Node,
	buildParserFromNodeFunc func(*buildContext[CT], *buildSubcontext[CT], Node) apc.Parser[CT, any]) apc.Parser[CT, []any] {

	parsers := make([]apc.Parser[CT, any], len(nodes))
	for i, rawNode := range nodes {
		parsers[i] = buildParserFromNodeFunc(buildCtx, subCtx, rawNode)
	}
	return apc.Seq(parsers...)
}

func buildAnyParserFromNodesCommon[CT any](buildCtx *buildContext[CT], subCtx *buildSubcontext[CT], nodes []Node,
	buildParserFromNodeFunc func(*buildContext[CT], *buildSubcontext[CT], Node) apc.Parser[CT, any]) apc.Parser[CT, any] {

	parsers := make([]apc.Parser[CT, any], len(nodes))
	for i, rawNode := range nodes {
		parsers[i] = buildParserFromNodeFunc(buildCtx, subCtx, rawNode)
	}
	return apc.Any(parsers...)
}

func buildParserFromNodeCommon[CT any](buildCtx *buildContext[CT], subCtx *buildSubcontext[CT], rawNode Node,
	buildParserForTypeFunc func(*buildContext[CT], reflect.Type) apc.Parser[CT, any],
	buildParserFromNodeFunc func(*buildContext[CT], *buildSubcontext[CT], Node) apc.Parser[CT, any]) apc.Parser[CT, any] {

	switch node := rawNode.(type) {
	case *inferNode:
		// TODO: if slice, use type of slice
		fieldName := subCtx.fieldNameFromCaptureIdx(node.InputIndex)
		field, ok := subCtx.resultStructType.FieldByName(fieldName)
		if !ok {
			panic(fmt.Sprintf("cannot infer parser: field '%v' not found in type '%v'", fieldName, subCtx.resultStructType.Name()))
		}
		if field.Type.Kind() == reflect.Slice {
			return buildParserForTypeFunc(buildCtx, field.Type.Elem())
		}
		return buildParserForTypeFunc(buildCtx, field.Type)
	case *captureNode:
		return apc.Map(
			buildParserFromNodeFunc(buildCtx, subCtx, node.Child),
			func(parseNode any) any {
				return captureResult{
					inputIndex: node.InputIndex,
					value:      parseNode,
				}
			},
		)
	case *seqNode:
		parser := buildSeqParserFromNodesCommon(buildCtx, subCtx, node.Children, buildParserFromNodeFunc)
		return apc.CastToAny(parser)
	case *rangeNode:
		childParser := buildParserFromNodeFunc(buildCtx, subCtx, node.Child)
		return apc.CastToAny(apc.Range(node.Range.min, node.Range.max, childParser))
	case *orNode:
		return buildAnyParserFromNodesCommon(buildCtx, subCtx, node.Children, buildParserFromNodeFunc)
	case *providedParserKeyNode:
		return buildCtx.mustGetProvidedParserByName(node.Name)
	case *maybeNode:
		return apc.CastToAny(apc.Maybe(buildParserFromNodeFunc(buildCtx, subCtx, node.Child)))
	case *lookNode:
		return apc.Look(buildParserFromNodeFunc(buildCtx, subCtx, node.Child))
	case *namedNode:
		return apc.Named(node.Name, buildParserFromNodeFunc(buildCtx, subCtx, node.Child))
	default:
		// To be handled in calling function
		return nil
	}
}

func valueSetFieldOrAppendKind(rawVal any, valKind reflect.Kind, fieldName string, field reflect.Value) {
	maybeAppendValToSliceTrueIfNot := func(val any) bool {
		if field.Kind() != reflect.Slice {
			return true
		}
		field.Set(reflect.Append(field, reflect.ValueOf(val)))
		return false
	}

	panicUnsettable := func(val any, exp string) {
		panic(fmt.Sprintf("cannot set field '%v' to value '%v' (type %T): cannot convert %v", fieldName, val, val, exp))
	}

	switch valKind {
	case reflect.String:
		switch val := rawVal.(type) {
		case string:
			if maybeAppendValToSliceTrueIfNot(val) {
				field.SetString(val)
			}
		default:
			strVal, ok := val.(string)
			if !ok {
				strVal = fmt.Sprintf("%v", val)
			}
			if maybeAppendValToSliceTrueIfNot(strVal) {
				field.SetString(strVal)
			}
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
			field.SetBool(val != reflect.Zero(reflect.TypeOf(val)).Interface())
		}
	case reflect.Pointer, reflect.Interface, reflect.Struct:
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
