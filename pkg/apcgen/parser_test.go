package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testOriginName = "testOrigin"
var emptyRootNode = &RootNode{Children: []Node{}}

func root1(child Node) *RootNode {
	return &RootNode{Children: []Node{child}}
}

func agg(children []Node) *AggregateNode {
	return &AggregateNode{Children: children}
}

func TestEmptyInput(t *testing.T) {
	node, err := parseFull(testOriginName, ``)
	assert.NoError(t, err)
	assert.Equal(t, emptyRootNode, node)

	node, err = parseFull(testOriginName, " \t  \t ")
	assert.NoError(t, err)
	assert.Equal(t, emptyRootNode, node)
}

func TestInfer(t *testing.T) {
	node, err := parseFull(testOriginName, `.`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		root1(&InferNode{
			InputIndex: 0,
		}),
		node)
}

func TestCaptureInfer(t *testing.T) {
	node, err := parseFull(testOriginName, `$.`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		root1(&CaptureNode{
			Child: &InferNode{
				InputIndex: 1,
			},
			InputIndex: 0,
		}),
		node)
}
