package apcgen

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

type captureResult struct {
	inputIndex int
	value      any
}

type buildContext[CT any] struct {
	options BuildOptions
	// Generated parser cache
	parserTypeParserCache map[reflect.Type]apc.Parser[CT, any]
	// Parser currently being generated but not yet done
	inProgressParserCache map[reflect.Type]*apc.Parser[CT, any]
	providedParserMap     map[string]apc.Parser[CT, any]
}

func newBuildContext[CT any](options BuildOptions, providedParsers map[string]apc.Parser[CT, any]) *buildContext[CT] {
	return &buildContext[CT]{
		options:               options,
		parserTypeParserCache: make(map[reflect.Type]apc.Parser[CT, any]),
		inProgressParserCache: make(map[reflect.Type]*apc.Parser[CT, any]),
		providedParserMap:     providedParsers,
	}
}

type buildSubcontext[CT any] struct {
	resultType                reflect.Type
	resultTypeElemName        string
	grammarText               string
	minCaptureIdxToFieldNames []keyValuePair[int, string]
}

func newBuildSubContextFromType[CT any](resultType reflect.Type) *buildSubcontext[CT] {
	var sb strings.Builder
	minCaptureIdxToFieldNames := make([]keyValuePair[int, string], 0)

	assertPointerToStructType(resultType)
	resultElemType := resultType.Elem()

	for i := 0; i < resultElemType.NumField(); i++ {
		field := resultElemType.Field(i)
		apcTag := field.Tag.Get("apc")
		if apcTag == "" {
			continue
		}
		sb.WriteString(apcTag)
		minCaptureIdxToFieldNames = append(minCaptureIdxToFieldNames, keyValuePair[int, string]{
			key:   sb.Len(),
			value: field.Name,
		})
	}

	return &buildSubcontext[CT]{
		resultType:                resultType,
		resultTypeElemName:        resultElemType.Name(),
		grammarText:               sb.String(),
		minCaptureIdxToFieldNames: minCaptureIdxToFieldNames,
	}
}

func (gc *buildContext[CT]) maybeGetCachedParserFromType(typ reflect.Type) apc.Parser[CT, any] {
	if parser, has := gc.parserTypeParserCache[typ]; has {
		return parser
	}
	return nil
}

func (gc *buildContext[CT]) maybeMakeRefParserFromType(typ reflect.Type) apc.Parser[CT, any] {
	if parserPtr, has := gc.inProgressParserCache[typ]; has {
		return apc.Ref(parserPtr)
	}
	return nil
}

func (gc *buildContext[CT]) mustGetProvidedParserByName(name string) apc.Parser[CT, any] {
	if parser, has := gc.providedParserMap[name]; has {
		return parser
	}
	panic(fmt.Sprintf("could not find a provided parser with name: %v", name))
}

func (gc *buildSubcontext[CT]) fieldNameFromCaptureIdx(idx int) string {
	// TODO: binary search or something
	for _, kvp := range gc.minCaptureIdxToFieldNames {
		if idx <= kvp.key {
			return kvp.value
		}
	}
	panic(fmt.Sprintf("could not find field name in struct '%v' relating to capture index %v known ranges were: %v",
		gc.resultTypeElemName, idx, gc.minCaptureIdxToFieldNames))
}

func assertPointerToStructType(resultType reflect.Type) {
	// Sanity checks for return type assumptions
	if resultType.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("the result type must be a pointer type; instead got: %v", resultType.Kind()))
	}

	resultTypeElem := resultType.Elem()
	if resultTypeElem.Kind() != reflect.Struct {
		panic(fmt.Sprintf("the result type must be a pointer to a struct; instead got pointer to: %v", resultTypeElem.Kind()))
	}
}
