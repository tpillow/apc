package apcgen

const (
	BuiltinMatchString string = "String"
	BuiltinMatchRegex  string = "Regex"
	BuiltinMatchToken  string = "Token"
)

type IntRange struct {
	Min int
	Max int
}

type Node interface{}

type RootNode struct {
	Children []Node
}

type MatchStringNode struct {
	Value string
}

type MatchRegexNode struct {
	Regex string
}

type MatchTokenNode struct {
	TokenTypeName string
}

type MatchTokenValueNode struct {
	TokenTypeName string
	StringValue   string
}

type InferNode struct {
	InputIndex int
}

type CaptureNode struct {
	Child      Node
	InputIndex int
}

type AggregateNode struct {
	Children []Node
}

type RangeNode struct {
	Range IntRange
	Child Node
}

type OrNode struct {
	Children []Node
}
