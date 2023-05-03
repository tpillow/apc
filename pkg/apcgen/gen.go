package apcgen

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/tpillow/apc/pkg/apc"
)

// TODO: maybe not global?
var runeParserTypeParserCache = make(map[reflect.Type]apc.Parser[rune, any])

type keyValuePair[KT, VT any] struct {
	key   KT
	value VT
}

type grammarDefinition struct {
	resultType                reflect.Type
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

func BuildRuneParser[RT any]() (apc.Parser[rune, RT], error) {
	rtType := reflect.TypeOf(new(RT))
	parser, err := buildRuneParserFromType(rtType)
	if err != nil {
		return nil, err
	}
	return apc.CastTo[rune, any, RT](parser), nil
}

func buildRuneParserFromType(resultType reflect.Type) (apc.Parser[rune, any], error) {
	if parser, has := runeParserTypeParserCache[resultType]; has {
		return parser, nil
	}

	grammarDef := getAllFieldsGrammarDef(resultType)
	node, err := parseFull(resultType.Name(), grammarDef.grammarText)
	if err != nil {
		return nil, err
	}
	parser, err := buildRuneParserFromRootNode(grammarDef, node)
	if err != nil {
		return nil, err
	}
	runeParserTypeParserCache[resultType] = parser
	return parser, nil
}

func buildRuneParserFromRootNode(grammarDef *grammarDefinition, node *RootNode) (apc.Parser[rune, any], error) {
	parser, err := buildRuneParserSequence(grammarDef, node.Children)
	if err != nil {
		return nil, err
	}

	return apc.Named(
		grammarDef.resultType.Name(),
		apc.Map(
			parser,
			func(parsedNodes []any) any {
				result := reflect.New(grammarDef.resultType)
				for i, childNode := range node.Children {
					switch node := childNode.(type) {
					case *CaptureNode:
						field := result.FieldByName(grammarDef.fieldNameFromCaptureIdx(node.InputIndex))
						field.SetPointer(unsafe.Pointer(&parsedNodes[i])) // TODO: ok?
					default:
					}
				}
				return result
			},
		),
	), nil
}

func buildRuneParserSequence(grammarDef *grammarDefinition, nodes []Node) (apc.Parser[rune, []any], error) {
	parsers := make([]apc.Parser[rune, any], len(nodes))
	for i, rawNode := range nodes {
		parser, err := buildRuneParserFromNode(grammarDef, rawNode)
		if err != nil {
			return nil, err
		}
		parsers[i] = parser
	}

	return apc.Seq(parsers...), nil
}

func buildRuneParserAny(grammarDef *grammarDefinition, nodes []Node) (apc.Parser[rune, any], error) {
	parsers := make([]apc.Parser[rune, any], len(nodes))
	for i, rawNode := range nodes {
		parser, err := buildRuneParserFromNode(grammarDef, rawNode)
		if err != nil {
			return nil, err
		}
		parsers[i] = parser
	}

	return apc.Any(parsers...), nil
}

func buildRuneParserFromNode(grammarDef *grammarDefinition, rawNode Node) (apc.Parser[rune, any], error) {
	switch node := rawNode.(type) {
	case *MatchStringNode:
		return apc.CastToAny(apc.ExactStr(node.Value)), nil
	case *MatchRegexNode:
		return apc.CastToAny(apc.Regex(node.Regex)), nil
	case *InferNode:
		fieldName := grammarDef.fieldNameFromCaptureIdx(node.InputIndex)
		field, ok := grammarDef.resultType.FieldByName(fieldName)
		if !ok {
			panic("could not find field with reflection")
		}
		return buildRuneParserFromType(field.Type)
	case *CaptureNode:
		return buildRuneParserFromNode(grammarDef, node.Child)
	case *AggregateNode:
		parser, err := buildRuneParserSequence(grammarDef, node.Children)
		if err != nil {
			return nil, err
		}
		return apc.CastToAny(parser), nil
	case *RangeNode:
		childParser, err := buildRuneParserFromNode(grammarDef, node.Child)
		if err != nil {
			return nil, err
		}
		return apc.CastToAny(apc.Range(node.Range.Min, node.Range.Max, childParser)), nil
	case *OrNode:
		parser, err := buildRuneParserAny(grammarDef, node.Children)
		if err != nil {
			return nil, err
		}
		return parser, nil
	default:
		panic("unknown node to process in buildRuneParserSequence")
	}
}

func getAllFieldsGrammarDef(resultType reflect.Type) *grammarDefinition {
	var sb strings.Builder
	minCaptureIdxToFieldNames := make([]keyValuePair[int, string], 0)

	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		apcTag := field.Tag.Get("apc")
		if apcTag == "" {
			continue
		}
		minCaptureIdxToFieldNames = append([]keyValuePair[int, string]{
			{
				key:   sb.Len(),
				value: field.Name,
			},
		}, minCaptureIdxToFieldNames...)
		sb.WriteString(apcTag)
	}

	return &grammarDefinition{
		resultType:                resultType,
		grammarText:               sb.String(),
		minCaptureIdxToFieldNames: minCaptureIdxToFieldNames,
	}
}
