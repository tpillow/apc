package apcgen

import (
	"fmt"

	"github.com/tpillow/apc/pkg/apc"
)

/*
parenExpr = '(' expr ')'
capturableValue = ( ident | '.' | '<str>' | string('') | regex('') | parenExpr )
valueMaybeCaptured = ( '$'? capturableValue )
value = valueMaybeCaptured endRangeSpecifier?

orExpr = value ('|' value)+
seqExpr = value value*

endRangeSpecifier = ('*'|'+'|'?'|{min,max})

expr = ( orExpr | seqExpr )
root = expr
*/

var (
	parserInitialized = false
	realExprParser    apc.Parser[rune, Node]
	exprParser        = apc.Ref(&realExprParser)

	parenExprParser = apc.Map(
		apc.Seq3(
			apc.Exact('('),
			exprParser,
			apc.Exact(')'),
		),
		func(node *apc.Seq3Node[rune, Node, rune]) Node {
			return node.Result2
		},
	)

	inferParser = apc.MapDetailed(
		apc.Exact('.'),
		func(_ rune, orgRange apc.OriginRange) (Node, error) {
			return &inferNode{
				InputIndex: orgRange.Start.ColNum,
			}, nil
		},
	)

	builtinStrLitParser = apc.Map(
		apc.SingleQuotedStringParser,
		func(node string) Node {
			return &matchStringNode{
				Value: node,
			}
		},
	)

	providedParserNameParser = apc.Map(
		apc.IdentifierParser,
		func(node string) Node {
			return &providedParserKeyNode{
				Name: node,
			}
		},
	)

	builtinMatchRegexParser = apc.Map(
		apc.Seq4(
			apc.ExactStr("regex"),
			apc.Exact('('),
			apc.SingleQuotedStringParser,
			apc.Exact(')'),
		),
		func(node *apc.Seq4Node[string, rune, string, rune]) Node {
			return &matchRegexNode{
				Regex: node.Result3,
			}
		},
	)

	builtinLookParser = apc.Map(
		apc.Seq4(
			apc.ExactStr("look"),
			apc.Exact('('),
			exprParser,
			apc.Exact(')'),
		),
		func(node *apc.Seq4Node[string, rune, Node, rune]) Node {
			return &lookNode{
				Child: node.Result3,
			}
		},
	)

	builtinFuncParser = apc.Any(
		builtinMatchRegexParser,
		builtinLookParser,
	)

	capturableValueParser = apc.Any(
		inferParser,
		builtinFuncParser,
		builtinStrLitParser,
		providedParserNameParser,
		parenExprParser,
	)

	valueMaybeCapturedParser = apc.Map(
		apc.Look(
			apc.Seq2(
				apc.MapDetailed(
					apc.Maybe(apc.Exact('$')),
					func(maybeNode apc.MaybeValue[rune], orgRange apc.OriginRange) (apc.MaybeValue[apc.Origin], error) {
						if !maybeNode.IsNil() {
							return apc.NewMaybeValue(orgRange.Start), nil
						}
						return apc.NewNilMaybeValue[apc.Origin](), nil
					},
				),
				capturableValueParser,
			),
		),
		func(node *apc.Seq2Node[apc.MaybeValue[apc.Origin], Node]) Node {
			if !node.Result1.IsNil() {
				return &captureNode{
					Child:      node.Result2,
					InputIndex: node.Result1.Value().ColNum,
				}
			}
			return node.Result2
		},
	)

	valueParser = apc.Map(
		apc.Seq2(
			valueMaybeCapturedParser,
			apc.Maybe(endRangeParser),
		),
		func(node *apc.Seq2Node[Node, apc.MaybeValue[intRange]]) Node {
			if node.Result2.IsNil() {
				return node.Result1
			}

			switch node.Result2.Value() {
			case maybeIntRange:
				rNode := &maybeNode{
					Child: node.Result1,
				}
				if capNode, ok := node.Result1.(*captureNode); ok {
					rNode.Child = capNode.Child
					return &captureNode{
						InputIndex: capNode.InputIndex,
						Child:      rNode,
					}
				}
				return rNode
			default:
				rNode := &rangeNode{
					Range: node.Result2.Value(),
					Child: node.Result1,
				}
				if capNode, ok := node.Result1.(*captureNode); ok {
					rNode.Child = capNode.Child
					return &captureNode{
						InputIndex: capNode.InputIndex,
						Child:      rNode,
					}
				}
				return rNode
			}
		},
	)

	seqExprParser = apc.Map(
		apc.OneOrMore(valueParser),
		func(nodes []Node) Node {
			if len(nodes) == 1 {
				return nodes[0]
			}
			return &seqNode{
				Children: nodes,
			}
		},
	)

	orExprParser = apc.Map(
		apc.Seq2(
			apc.Look(
				apc.Seq2(
					valueParser,
					apc.Exact('|'),
				),
			),
			apc.OneOrMoreSeparated(
				valueParser,
				apc.Exact('|'),
			),
		),
		func(node *apc.Seq2Node[*apc.Seq2Node[Node, rune], []Node]) Node {
			children := []Node{node.Result1.Result1}
			children = append(children, node.Result2...)
			return &orNode{
				Children: children,
			}
		},
	)

	endRangeParser = apc.Any(
		apc.Map(
			apc.Regex("[\\*\\+\\?]"),
			func(node string) intRange {
				switch node {
				case "*":
					return intRange{
						min: 0,
						max: -1,
					}
				case "+":
					return intRange{
						min: 1,
						max: -1,
					}
				case "?":
					return maybeIntRange
				default:
					panic("unreachable in endRangeParser")
				}
			},
		),
		apc.MapDetailed(
			apc.Seq5(
				apc.Exact('{'),
				apc.IntParser,
				apc.Exact(','),
				apc.IntParser,
				apc.Exact('}'),
			),
			func(node *apc.Seq5Node[rune, int64, rune, int64, rune], _ apc.OriginRange) (intRange, error) {
				ir := intRange{
					min: int(node.Result2),
					max: int(node.Result4),
				}
				if ir.min < 0 {
					return intRange{},
						fmt.Errorf("invalid int range {%v, %v}: min value must be >= 0", ir.min, ir.max)
				}
				if ir.max < 0 && ir.max != -1 {
					return intRange{},
						fmt.Errorf("invalid int range {%v, %v}: max value must be >= 0 (or be == -1 for no limit)", ir.min, ir.max)
				}
				return ir, nil
			},
		),
	)

	rootParser = apc.Skip(
		// TODO: not allowing newlines, so we can track position via the column
		// of an Origin to map captures to appropriate fields. Can probably
		// come up with a work-around at some point.
		apc.CastToAny(apc.Regex("( |\\t)+")),
		apc.Map(
			exprParser,
			func(node Node) *rootNode {
				return &rootNode{
					Child: node,
				}
			},
		),
	)
)

func maybeInitParser() {
	if parserInitialized {
		return
	}
	parserInitialized = true

	realExprParser = apc.Any(
		orExprParser,
		seqExprParser,
	)
}

func parseFull(originName string, input string, debugParsers bool) (*rootNode, error) {
	maybeInitParser()
	ctx := apc.NewStringContext(originName, input)
	ctx.DebugParsers = debugParsers
	return apc.Parse[rune](ctx, rootParser, apc.DefaultParseConfig)
}
