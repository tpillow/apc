package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/apc/pkg/apc"
)

var testOriginName = "testOrigin"
var agg0 = &AggregateNode{Children: []Node{}}

func agg1(child Node) *AggregateNode {
	return &AggregateNode{Children: []Node{child}}
}

func agg(children []Node) *AggregateNode {
	return &AggregateNode{Children: children}
}

func orgCol(col int) apc.Origin {
	return apc.Origin{
		Name:    testOriginName,
		LineNum: 1,
		ColNum:  col,
	}
}

func TestEmptyInput(t *testing.T) {
	node, err := parseFull(testOriginName, ``)
	assert.NoError(t, err)
	assert.Equal(t, agg0, node)

	node, err = parseFull(testOriginName, " \t  \t ")
	assert.NoError(t, err)
	assert.Equal(t, agg0, node)
}

func TestInfer(t *testing.T) {
	node, err := parseFull(testOriginName, `.`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		agg1(&InferNode{}),
		node)
}

func TestCaptureInfer(t *testing.T) {
	node, err := parseFull(testOriginName, `$.`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		agg1(&CaptureNode{
			Child:  &InferNode{},
			Origin: orgCol(1),
		}),
		node)
}
