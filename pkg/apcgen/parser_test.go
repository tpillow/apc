package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testOriginName = "testOrigin"

func root1(child Node) *RootNode {
	return &RootNode{Child: child}
}

func TestEmptyInput(t *testing.T) {
	_, err := parseFull(testOriginName, ``)
	assert.Error(t, err)
	_, err = parseFull(testOriginName, " \t  \t ")
	assert.Error(t, err)
}

func TestInfer(t *testing.T) {
	node, err := parseFull(testOriginName, `.`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		root1(&InferNode{
			InputIndex: 1,
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
				InputIndex: 2,
			},
			InputIndex: 1,
		}),
		node)
}
