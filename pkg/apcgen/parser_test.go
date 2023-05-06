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

func TestGeneric1(t *testing.T) {
	node, err := parseFull(
		testOriginName,
		`'Entry' '{' $StrParser $regex('[0-9]+') $.? $(.*) '}'`)
	assert.NoError(t, err)
	assert.Equal(
		t,
		root1(
			&SeqNode{
				Children: []Node{
					&MatchStringNode{Value: "Entry"},
					&MatchStringNode{Value: "{"},
					&CaptureNode{
						InputIndex: 13,
						Child:      &ProvidedParserKeyNode{Name: "StrParser"},
					},
					&CaptureNode{
						InputIndex: 24,
						Child:      &MatchRegexNode{Regex: "[0-9]+"},
					},
					&CaptureNode{
						InputIndex: 41,
						Child: &MaybeNode{
							Child: &InferNode{
								InputIndex: 42,
							},
						},
					},
					&CaptureNode{
						InputIndex: 44,
						Child: &RangeNode{
							Range: IntRange{Min: 0, Max: -1},
							Child: &InferNode{
								InputIndex: 47,
							},
						},
					},
					&MatchStringNode{Value: "}"},
				},
			},
		),
		node)
}
