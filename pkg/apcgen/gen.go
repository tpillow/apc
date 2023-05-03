package apcgen

import (
	"reflect"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

type keyValuePair[KT, VT any] struct {
	key   KT
	value VT
}

type grammarDefinition struct {
	grammarText               string
	minCaptureIdxToFieldNames []keyValuePair[int, string]
}

func (gd *grammarDefinition) fieldNameFromCaptureIdx(idx int) string {
	// TODO: binary search or something
	for _, kvp := range gd.minCaptureIdxToFieldNames {
		if idx <= kvp.key {
			return kvp.value
		}
	}
	panic("could not find field name in struct relating to capture index")
}

func BuildRuneParser[RT any](node *RootNode) (apc.Parser[rune, RT], error) {
	rtType := reflect.TypeOf(new(RT))

	grammarDef := getAllFieldsGrammarDef[RT]()
	node, err := parseFull(rtType.Name(), grammarDef.grammarText)
	if err != nil {
		return nil, err
	}

	parser, err := buildRuneParserFromRootNode(rtType, grammarDef, node)
	if err != nil {
		return nil, err
	}
	return apc.CastTo[rune, any, RT](parser), nil
}

func buildRuneParserFromRootNode(resultType reflect.Type, grammarDef *grammarDefinition, node *RootNode) (apc.Parser[rune, any], error) {
	return nil, nil
}

func buildRuneParserSequence(nodes []Node) (apc.Parser[rune, []any], error) {
	return nil, nil
}

func getAllFieldsGrammarDef[T any]() *grammarDefinition {
	var sb strings.Builder
	minCaptureIdxToFieldNames := make([]keyValuePair[int, string], 0)

	resultType := reflect.TypeOf(new(T))
	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		apcTag := field.Tag.Get("apc")
		if apcTag == "" {
			continue
		}
		minCaptureIdxToFieldNames = append([]keyValuePair[int, string]{
			keyValuePair[int, string]{
				key:   sb.Len(),
				value: field.Name,
			},
		}, minCaptureIdxToFieldNames...)
		sb.WriteString(apcTag)
	}

	return &grammarDefinition{
		grammarText:               sb.String(),
		minCaptureIdxToFieldNames: minCaptureIdxToFieldNames,
	}
}
