package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
	opAddParser = apc.Bind(apc.ExactStr(string(OpAdd)), OpAdd)
	opSubParser = apc.Bind(apc.ExactStr(string(OpSub)), OpSub)
	opMulParser = apc.Bind(apc.ExactStr(string(OpMul)), OpMul)
	opDivParser = apc.Bind(apc.ExactStr(string(OpDiv)), OpDiv)

	factorParser    apc.Parser[string, Executable]
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

	maybeExprParser = apc.Maybe("expr", exprParser)
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
				apc.ExactStr("("),
				exprParser,
				apc.ExactStr(")")),
			func(node *apc.Seq3Node[string, Executable, string]) Executable {
				return node.Result2
			}))
}

func executeInput(input string) {
	ctx := apc.NewStringContext("<user_input>", input)
	ctx.AddSkipParser(apc.MapToAny(apc.WhitespaceParser))

	node, err := apc.Parse[string](ctx, maybeExprParser, apc.DefaultParseConfig)
	if err != nil {
		fmt.Printf("Error parsing input: %v\n", err)
		return
	}
	if node == nil {
		return
	}

	fmt.Printf("Parse Node:       %v\n", node)
	fmt.Printf("Result:           %v\n\n", node.Execute())
}

func main() {
	initParser()
	inputPrompt := "Input expression: "

	input := "1 + 2 * (4 - 3)"
	fmt.Print("Welcome to the APC calculator. Type 'q' to quit.\n")
	fmt.Printf("Here, I'll do one first:\n\n%v%v\n", inputPrompt, input)
	executeInput(input)

	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(inputPrompt)
		input, err := stdin.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		if strings.TrimSpace(input) == "q" {
			break
		}
		executeInput(input)
	}

	fmt.Printf("Goodbye!\n")
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
