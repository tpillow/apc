package main

import (
	"bufio"
	"fmt"
	"math"
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
	OpExp Operator = "^"
)

var (
	opAddParser = apc.Bind(apc.ExactStr(string(OpAdd)), OpAdd)
	opSubParser = apc.Bind(apc.ExactStr(string(OpSub)), OpSub)
	opMulParser = apc.Bind(apc.ExactStr(string(OpMul)), OpMul)
	opDivParser = apc.Bind(apc.ExactStr(string(OpDiv)), OpDiv)
	opExpParser = apc.Bind(apc.ExactStr(string(OpExp)), OpExp)

	factorParser    apc.Parser[rune, Executable]
	factorParserRef = apc.Ref(&factorParser)

	exponentTermParser = apc.Named("exponential term",
		apc.Map(
			apc.Seq2(
				factorParserRef,
				apc.ZeroOrMore(
					apc.Seq2(
						opExpParser,
						factorParserRef))),
			func(node *apc.Seq2Node[Executable, []*apc.Seq2Node[Operator, Executable]], _ apc.Origin) Executable {
				left := node.Result1
				for _, seqRes := range node.Result2 {
					left = BinOpNode{
						Operator: seqRes.Result1,
						Left:     left,
						Right:    seqRes.Result2,
					}
				}
				return left
			}))

	termParser = apc.Named("term",
		apc.Map(
			apc.Seq2(
				exponentTermParser,
				apc.ZeroOrMore(
					apc.Seq2(
						apc.Any(opMulParser, opDivParser),
						exponentTermParser))),
			func(node *apc.Seq2Node[Executable, []*apc.Seq2Node[Operator, Executable]], _ apc.Origin) Executable {
				left := node.Result1
				for _, seqRes := range node.Result2 {
					left = BinOpNode{
						Operator: seqRes.Result1,
						Left:     left,
						Right:    seqRes.Result2,
					}
				}
				return left
			}))

	exprParser = apc.Named("expression",
		apc.Map(
			apc.Seq2(
				termParser,
				apc.ZeroOrMore(
					apc.Seq2(
						apc.Any(opAddParser, opSubParser),
						termParser))),
			func(node *apc.Seq2Node[Executable, []*apc.Seq2Node[Operator, Executable]], _ apc.Origin) Executable {
				left := node.Result1
				for _, seqRes := range node.Result2 {
					left = BinOpNode{
						Operator: seqRes.Result1,
						Left:     left,
						Right:    seqRes.Result2,
					}
				}
				return left
			}))

	maybeExprParser = apc.Maybe(exprParser)
)

func initParser() {
	factorParser = apc.Named("factor",
		apc.Any(
			apc.Map(
				apc.FloatParser,
				func(node float64, _ apc.Origin) Executable {
					return ValueNode{Value: node}
				}),
			apc.Map(
				apc.Seq3(
					apc.ExactStr("("),
					exprParser,
					apc.ExactStr(")")),
				func(node *apc.Seq3Node[string, Executable, string], _ apc.Origin) Executable {
					return node.Result2
				})))
}

func executeInput(input string) {
	ctx := apc.NewStringContext("<user_input>", input)
	ctx.AddSkipParser(apc.CastToAny(apc.WhitespaceParser))

	node, err := apc.Parse[rune](ctx, maybeExprParser, apc.DefaultParseConfig)
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

	input := "1 + 2 * 3 ^ (5 - 3)"
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
	case OpExp:
		return math.Pow(n.Left.Execute(), n.Right.Execute())
	}
	panic("unknown operator")
}
