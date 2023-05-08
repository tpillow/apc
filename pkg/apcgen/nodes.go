package apcgen

type Node interface{}

type rootNode struct {
	Child Node
}

type providedParserKeyNode struct {
	Name string
}

type matchStringNode struct {
	Value string
}

type matchRegexNode struct {
	Regex string
}

type inferNode struct {
	InputIndex int
}

type captureNode struct {
	Child      Node
	InputIndex int
}

type seqNode struct {
	Children []Node
}

type orNode struct {
	Children []Node
}
type rangeNode struct {
	Range intRange
	Child Node
}

type maybeNode struct {
	Child Node
}

type lookNode struct {
	Child Node
}
