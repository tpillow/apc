package apcgen

import "github.com/tpillow/apc/pkg/apc"

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
	Origin apc.Origin
}

type CaptureNode struct {
	Child  Node
	Origin apc.Origin
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
