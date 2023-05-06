package apcgen

import "fmt"

type Node interface{}

type RootNode struct {
	Child Node
}

func (node *RootNode) String() string {
	return fmt.Sprintf("<Root Child=%v>", node.Child)
}

type ProvidedParserKeyNode struct {
	Name string
}

func (node *ProvidedParserKeyNode) String() string {
	return fmt.Sprintf("<ProvidedParserKeyNode Name=%v>", node.Name)
}

type MatchStringNode struct {
	Value string
}

func (node *MatchStringNode) String() string {
	return fmt.Sprintf("<MatchStringNode Value=%v>", node.Value)
}

type MatchRegexNode struct {
	Regex string
}

func (node *MatchRegexNode) String() string {
	return fmt.Sprintf("<MatchRegexNode Regex=%v>", node.Regex)
}

type InferNode struct {
	InputIndex int
}

func (node *InferNode) String() string {
	return fmt.Sprintf("<InferNode InputIndex=%v>", node.InputIndex)
}

type CaptureNode struct {
	Child      Node
	InputIndex int
}

func (node *CaptureNode) String() string {
	return fmt.Sprintf("<CaptureNode InputIndex=%v Child=%v>", node.InputIndex, node.Child)
}

type SeqNode struct {
	Children []Node
}

func (node *SeqNode) String() string {
	return fmt.Sprintf("<SeqNode Children=%v>", node.Children)
}

type OrNode struct {
	Children []Node
}

func (node *OrNode) String() string {
	return fmt.Sprintf("<OrNode Children=%v>", node.Children)
}

type RangeNode struct {
	Range IntRange
	Child Node
}

func (node *RangeNode) String() string {
	return fmt.Sprintf("<RangeNode Range=%v Child=%v>", node.Range, node.Child)
}

type MaybeNode struct {
	Child Node
}

func (node *MaybeNode) String() string {
	return fmt.Sprintf("<MaybeNode Child=%v>", node.Child)
}
