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

	factorParser = apc.FloatParser

	termParser = apc.Map(
		apc.Seq2("term",
			factorParser,
			apc.Maybe("",
				apc.Seq2("",
					apc.OneOf("", opMulParser, opDivParser),
					factorParser))),
		func(node *apc.Seq2Node[float64, *apc.Seq2Node[Operator, float64]]) Executable {
			if node.Result2 == nil {
				return ValueNode{
					Value: node.Result1,
				}
			}
			return BinOpNode{
				Operator: node.Result2.Result1,
				Left: ValueNode{
					Value: node.Result1,
				},
				Right: ValueNode{
					Value: node.Result2.Result2,
				},
			}
		})

	exprParser = apc.Map(
		apc.Seq2("expr",
			termParser,
			apc.Maybe("",
				apc.Seq2("",
					apc.OneOf("", opAddParser, opSubParser),
					termParser))),
		func(node *apc.Seq2Node[Executable, *apc.Seq2Node[Operator, Executable]]) Executable {
			if node.Result2 == nil {
				return node.Result1
			}
			return BinOpNode{
				Operator: node.Result2.Result1,
				Left:     node.Result1,
				Right:    node.Result2.Result2,
			}
		})
)

func main() {
	input := "11 * 22 + 33 * 44"
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
