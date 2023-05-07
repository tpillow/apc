package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testOriginName = "testOrigin"

func root1(child Node) *rootNode {
	return &rootNode{Child: child}
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
		root1(&inferNode{
			InputIndex: 1,
		}),
		node)
}

func TestCaptureInfer(t *testing.T) {
	node, err := parseFull(testOriginName, `$.`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		root1(&captureNode{
			Child: &inferNode{
				InputIndex: 2,
			},
			InputIndex: 1,
		}),
		node)
}

func TestCaptureInferRange(t *testing.T) {
	node, err := parseFull(testOriginName, `$.{1,3}`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		root1(&captureNode{
			Child: &rangeNode{
				Range: intRange{min: 1, max: 3},
				Child: &inferNode{
					InputIndex: 2,
				},
			},
			InputIndex: 1,
		}),
		node)
}

func TestGeneric1(t *testing.T) {
	node, err := parseFull(
		testOriginName,
		`'Entry' '{' $StrParser $regex('[0-9]+') $.? $(.*) '}'`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		root1(
			&seqNode{
				Children: []Node{
					&matchStringNode{Value: "Entry"},
					&matchStringNode{Value: "{"},
					&captureNode{
						InputIndex: 13,
						Child:      &providedParserKeyNode{Name: "StrParser"},
					},
					&captureNode{
						InputIndex: 24,
						Child:      &matchRegexNode{Regex: "[0-9]+"},
					},
					&captureNode{
						InputIndex: 41,
						Child: &maybeNode{
							Child: &inferNode{
								InputIndex: 42,
							},
						},
					},
					&captureNode{
						InputIndex: 44,
						Child: &rangeNode{
							Range: intRange{min: 0, max: -1},
							Child: &inferNode{
								InputIndex: 47,
							},
						},
					},
					&matchStringNode{Value: "}"},
				},
			},
		),
		node)
}
