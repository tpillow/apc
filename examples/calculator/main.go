package main

import (
	"fmt"

	"github.com/tpillow/apc/pkg/apc"
)

type Executable interface {
	Execute() float64
}

type Operator string

const (
	OpAdd Operator = "+"
	OpSub Operator = "-"
	OpMul Operator = "*"
	OpDiv Operator = "/"
)

var (
	opAddParser = apc.Bind(apc.Exact(string(OpAdd)), OpAdd)
	opSubParser = apc.Bind(apc.Exact(string(OpSub)), OpSub)
	opMulParser = apc.Bind(apc.Exact(string(OpMul)), OpMul)
	opDivParser = apc.Bind(apc.Exact(string(OpDiv)), OpDiv)

	factorParser    apc.Parser[Executable]
	factorParserRef = apc.Ref(&factorParser)

	termParser = apc.Map(
		apc.Seq2("term",
			factorParserRef,
			apc.ZeroOrMore("",
				apc.Seq2("",
					apc.OneOf("", opMulParser, opDivParser),
					factorParserRef))),
		func(node *apc.Seq2Node[Executable, []*apc.Seq2Node[Operator, Executable]]) Executable {
			left := node.Result1
			for _, seqRes := range node.Result2 {
				left = BinOpNode{
					Operator: seqRes.Result1,
					Left:     left,
					Right:    seqRes.Result2,
				}
			}
			return left
		})

	exprParser = apc.Map(
		apc.Seq2("expr",
			termParser,
			apc.ZeroOrMore("",
				apc.Seq2("",
					apc.OneOf("", opAddParser, opSubParser),
					termParser))),
		func(node *apc.Seq2Node[Executable, []*apc.Seq2Node[Operator, Executable]]) Executable {
			left := node.Result1
			for _, seqRes := range node.Result2 {
				left = BinOpNode{
					Operator: seqRes.Result1,
					Left:     left,
					Right:    seqRes.Result2,
				}
			}
			return left
		})
)

func initParser() {
	factorParser = apc.OneOf("factor",
		apc.Map(
			apc.FloatParser,
			func(node float64) Executable {
				return ValueNode{Value: node}
			}),
		apc.Map(
			apc.Seq3("",
				apc.Exact("("),
				exprParser,
				apc.Exact(")")),
			func(node *apc.Seq3Node[string, Executable, string]) Executable {
				return node.Result2
			}))
}

func main() {
	initParser()

	input := "11 * (22 + 33) * 44"
	ctx := apc.NewStringContext("<string>", input)
	ctx.AddSkipParser(apc.MapToAny(apc.WhitespaceParser))

	node, err := apc.Parse(ctx, exprParser, apc.DefaultParseConfig)

	fmt.Printf("Input: %v\n", input)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Result Node: %v\n", node)
	fmt.Printf("Result: %v\n", node.Execute())
}

type ValueNode struct {
	Value float64
}

func (n ValueNode) Execute() float64 {
	return n.Value
}

type BinOpNode struct {
	Operator Operator
	Left     Executable
	Right    Executable
}

func (n BinOpNode) Execute() float64 {
	switch n.Operator {
	case OpAdd:
		return n.Left.Execute() + n.Right.Execute()
	case OpSub:
		return n.Left.Execute() - n.Right.Execute()
	case OpMul:
		return n.Left.Execute() * n.Right.Execute()
	case OpDiv:
		return n.Left.Execute() / n.Right.Execute()
	}
	panic("unknown operator")
}
