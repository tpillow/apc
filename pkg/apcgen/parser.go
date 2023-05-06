package apcgen

import (
	"github.com/tpillow/apc/pkg/apc"
)

/*
parenExpr = '(' expr ')'
capturableValue = ( ident | '.' | '<str>' | string('') | regex('') | parenExpr )
valueMaybeCaptured = ( '$'? capturableValue )
value = valueMaybeCaptured endRangeSpecifier?

orExpr = value ('|' value)+
seqExpr = value value*

endRangeSpecifier = ('*'|'+'|'?')

expr = ( orExpr | seqExpr )
root = expr
*/

var (
	realExprParser apc.Parser[rune, Node]
	exprParser     = apc.Ref(&realExprParser)

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
			return &InferNode{
				InputIndex: orgRange.Start.ColNum,
			}, nil
		},
	)

	builtinStrLitParser = apc.Map(
		apc.SingleQuotedStringParser,
		func(node string) Node {
			return &MatchStringNode{
				Value: node,
			}
		},
	)

	providedParserNameParser = apc.Map(
		apc.IdentifierParser,
		func(node string) Node {
			return &ProvidedParserKeyNode{
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
			return &MatchRegexNode{
				Regex: node.Result3,
			}
		},
	)

	builtinFuncParser = builtinMatchRegexParser

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
					apc.Maybe(apc.ExactStr("$")),
					func(maybeNode apc.MaybeValue[string], orgRange apc.OriginRange) (apc.MaybeValue[apc.Origin], error) {
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
				return &CaptureNode{
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
			endRangeParser,
		),
		func(node *apc.Seq2Node[Node, apc.MaybeValue[IntRange]]) Node {
			if node.Result2.IsNil() {
				return node.Result1
			}

			switch node.Result2.Value() {
			case maybeIntRange:
				rNode := &MaybeNode{
					Child: node.Result1,
				}
				if capNode, ok := node.Result1.(*CaptureNode); ok {
					rNode.Child = capNode.Child
					return &CaptureNode{
						InputIndex: capNode.InputIndex,
						Child:      rNode,
					}
				}
				return rNode
			default:
				rNode := &RangeNode{
					Range: node.Result2.Value(),
					Child: node.Result1,
				}
				if capNode, ok := node.Result1.(*CaptureNode); ok {
					rNode.Child = capNode.Child
					return &CaptureNode{
						InputIndex: capNode.InputIndex,
						Child:      rNode,
					}
				}
				return rNode
			}
		},
	)

	seqExprParser = apc.Map(
		apc.Seq2(
			valueParser,
			apc.ZeroOrMore(valueParser),
		),
		func(node *apc.Seq2Node[Node, []Node]) Node {
			if len(node.Result2) == 0 {
				return node.Result1
			}

			children := []Node{node.Result1}
			for _, child := range node.Result2 {
				children = append(children, child)
			}
			return &SeqNode{
				Children: children,
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
			for _, child := range node.Result2 {
				children = append(children, child)
			}
			return &OrNode{
				Children: children,
			}
		},
	)

	endRangeParser = apc.Map(
		// TODO: any range specifier {min, max}
		apc.Regex("[\\*\\+\\?]?"),
		func(node string) apc.MaybeValue[IntRange] {
			if node == "" {
				return apc.NewNilMaybeValue[IntRange]()
			}

			switch node {
			case "*":
				return apc.NewMaybeValue(IntRange{
					Min: 0,
					Max: -1,
				})
			case "+":
				return apc.NewMaybeValue(IntRange{
					Min: 1,
					Max: -1,
				})
			case "?":
				return apc.NewMaybeValue(maybeIntRange)
			default:
				panic("unreachable in endRangeParser")
			}
		},
	)

	rootParser = apc.Skip(
		// TODO: not allowing newlines, so we can track position via the column
		// of an Origin to map captures to appropriate fields. Can probably
		// come up with a work-around at some point.
		apc.CastToAny(apc.Regex("[ \t]+")),
		apc.Map(
			exprParser,
			func(node Node) *RootNode {
				return &RootNode{
					Child: node,
				}
			},
		),
	)
)

func initParser() {
	realExprParser = apc.Any(
		orExprParser,
		seqExprParser,
	)
}

func parseFull(originName string, input string) (*RootNode, error) {
	initParser()
	ctx := apc.NewStringContext(originName, input)
	return apc.Parse[rune](ctx, rootParser, apc.DefaultParseConfig)
}
