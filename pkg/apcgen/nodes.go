package apcgen

import "fmt"

type Node interface{}

type rootNode struct {
	Child Node
}

func (node *rootNode) String() string {
	return fmt.Sprintf("<Root Child=%v>", node.Child)
}

type providedParserKeyNode struct {
	Name string
}

func (node *providedParserKeyNode) String() string {
	return fmt.Sprintf("<ProvidedParserKeyNode Name=%v>", node.Name)
}

type matchStringNode struct {
	Value string
}

func (node *matchStringNode) String() string {
	return fmt.Sprintf("<MatchStringNode Value=%v>", node.Value)
}

type matchRegexNode struct {
	Regex string
}

func (node *matchRegexNode) String() string {
	return fmt.Sprintf("<MatchRegexNode Regex=%v>", node.Regex)
}

type inferNode struct {
	InputIndex int
}

func (node *inferNode) String() string {
	return fmt.Sprintf("<InferNode InputIndex=%v>", node.InputIndex)
}

type captureNode struct {
	Child      Node
	InputIndex int
}

func (node *captureNode) String() string {
	return fmt.Sprintf("<CaptureNode InputIndex=%v Child=%v>", node.InputIndex, node.Child)
}

type seqNode struct {
	Children []Node
}

func (node *seqNode) String() string {
	return fmt.Sprintf("<SeqNode Children=%v>", node.Children)
}

type orNode struct {
	Children []Node
}

func (node *orNode) String() string {
	return fmt.Sprintf("<OrNode Children=%v>", node.Children)
}

type rangeNode struct {
	Range intRange
	Child Node
}

func (node *rangeNode) String() string {
	return fmt.Sprintf("<RangeNode Range=%v Child=%v>", node.Range, node.Child)
}

type maybeNode struct {
	Child Node
}

func (node *maybeNode) String() string {
	return fmt.Sprintf("<MaybeNode Child=%v>", node.Child)
}
