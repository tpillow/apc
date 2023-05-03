package apcgen

type Node interface{}

type RootNode struct {
	Child Node
}

type ProvidedParserKeyNode struct {
	Name string
}

type MatchStringNode struct {
	Value string
}

type MatchRegexNode struct {
	Regex string
}

type InferNode struct {
	InputIndex int
}

type CaptureNode struct {
	Child      Node
	InputIndex int
}

type SeqNode struct {
	Children []Node
}

type RangeNode struct {
	Range IntRange
	Child Node
}

type OrNode struct {
	Children []Node
}
