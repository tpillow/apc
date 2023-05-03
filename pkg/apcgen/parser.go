package apcgen

import (
	"fmt"

	"github.com/tpillow/apc/pkg/apc"
)

var (
	parserInitialized = false

	exprParser    apc.Parser[rune, Node]
	exprParserRef = apc.Ref(&exprParser)

	strParamListParser = apc.Named("string parameter list",
		apc.ZeroOrMoreSeparated(
			apc.SingleQuotedStringParser,
			apc.ExactStr(","),
		))

	matchStrParser = apc.Named("implicit match string",
		apc.Map(
			apc.SingleQuotedStringParser,
			func(node string) Node {
				return &MatchStringNode{
					Value: node,
				}
			},
		))

	builtinCallParser = apc.Named("builtin call",
		apc.MapDetailed(
			apc.Seq4(
				apc.Regex(fmt.Sprintf("%v|%v|%v",
					BuiltinMatchString, BuiltinMatchRegex, BuiltinMatchToken)),
				apc.ExactStr("("),
				strParamListParser,
				apc.ExactStr(")"),
			),
			func(node *apc.Seq4Node[string, string, []string, string], _ apc.OriginRange) (Node, error) {
				switch node.Result1 {
				case BuiltinMatchString:
					switch len(node.Result3) {
					case 1:
						return &MatchStringNode{
							Value: node.Result3[0],
						}, nil
					default:
						return nil, fmt.Errorf("Builtin function '%v' must have only 1 parameter, but got: %v", node.Result1, node.Result3)
					}
				case BuiltinMatchRegex:
					switch len(node.Result3) {
					case 1:
						return &MatchRegexNode{
							Regex: node.Result3[0],
						}, nil
					default:
						return nil, fmt.Errorf("Builtin function '%v' must have only 1 parameter, but got: %v", node.Result1, node.Result3)
					}
				case BuiltinMatchToken:
					switch len(node.Result3) {
					case 1:
						return &MatchTokenNode{
							TokenTypeName: node.Result3[0],
						}, nil
					case 2:
						return &MatchTokenValueNode{
							TokenTypeName: node.Result3[0],
							StringValue:   node.Result3[1],
						}, nil
					default:
						return nil, fmt.Errorf("builtin function '%v' must have only 1 or 2 parameters, but got: %v", node.Result1, node.Result3)
					}
				default:
					panic("unreachable")
				}
			},
		))

	rawValueParser = apc.Named("value",
		apc.Map(
			apc.OneOrMoreSeparated(
				apc.Any(
					apc.MapDetailed(
						apc.Exact('.'),
						func(node rune, orgRange apc.OriginRange) (Node, error) {
							return &InferNode{
								InputIndex: orgRange.Start.ColNum,
							}, nil
						},
					),
					builtinCallParser,
					matchStrParser,
				),
				apc.ExactStr("|"),
			),
			func(nodes []Node) Node {
				if len(nodes) == 1 {
					return nodes[0]
				}
				return &AggregateNode{
					Children: nodes,
				}
			},
		))

	valueParser = apc.Named("maybe-capture value",
		apc.Map(
			apc.Look(
				apc.Seq2(
					apc.MapDetailed(
						apc.Maybe(apc.ExactStr("$")),
						func(maybeNode apc.MaybeValue[string], orgRange apc.OriginRange) (apc.MaybeValue[apc.Origin], error) {
							if !maybeNode.IsNil() {
								return apc.NewMaybeValue(orgRange.Start), nil
							}
							return apc.NewNilMaybeValue[apc.Origin](), nil
						},
					),
					rawValueParser,
				),
			),
			func(node *apc.Seq2Node[apc.MaybeValue[apc.Origin], Node]) Node {
				if !node.Result1.IsNil() {
					return &CaptureNode{
						Child:      node.Result2,
						InputIndex: node.Result1.Value().ColNum,
					}
				}
				return node.Result2
			},
		))

	parenExprParser = apc.Named("parenthesized expression",
		apc.Map(
			apc.Seq3(
				apc.ExactStr("("),
				exprParserRef,
				apc.ExactStr(")"),
			),
			func(node *apc.Seq3Node[string, Node, string]) Node {
				return node.Result2
			},
		))

	endRangeModifierParser = apc.Named("end-range modifier",
		apc.Map(
			// TODO: generic ranges {min, max}
			apc.Regex("\\*|\\+|\\?"),
			func(node string) IntRange {
				valMap := map[string]IntRange{
					"*": {0, -1},
					"+": {1, -1},
					"?": {0, 1},
				}
				if val, ok := valMap[node]; ok {
					return val
				}
				panic("unreachable")
			},
		))

	rootParser = apc.Skip(
		// TODO: not allowing newlines, so we can track position via the column
		// of an Origin to map captures to appropriate fields. Can probably
		// come up with a work-around at some point.
		apc.CastToAny(apc.Regex("[ \t]+")),
		apc.Map(
			apc.ZeroOrMore(exprParserRef),
			func(nodes []Node) *RootNode {
				return &RootNode{
					Children: nodes,
				}
			},
		),
	)
)

func initParser() {
	if parserInitialized {
		return
	}
	parserInitialized = true

	exprParser = apc.Named("expression",
		apc.Map(
			apc.Seq2(
				apc.Any(
					parenExprParser,
					valueParser,
				),
				apc.Maybe(
					endRangeModifierParser,
				),
			),
			func(node *apc.Seq2Node[Node, apc.MaybeValue[IntRange]]) Node {
				if !node.Result2.IsNil() {
					return &RangeNode{
						Range: node.Result2.Value(),
						Child: node.Result1,
					}
				}
				return node.Result1
			},
		))
}

func parseFull(originName string, input string) (*RootNode, error) {
	initParser()
	ctx := apc.NewStringContext(originName, input)
	return apc.Parse[rune](ctx, rootParser, apc.DefaultParseConfig)
}
