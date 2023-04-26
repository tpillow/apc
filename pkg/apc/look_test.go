package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLook(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "abcabdabeabcabfabd")

	parser := Any(
		Look(Map(Seq3(ExactStr("a"), ExactStr("b"), ExactStr("c")),
			func(node *Seq3Node[string, string, string], _ Origin) string {
				return "c"
			})),
		Look(Map(Seq3(ExactStr("a"), ExactStr("b"), ExactStr("d")),
			func(node *Seq3Node[string, string, string], _ Origin) string {
				return "d"
			})),
		Any(
			Look(Map(Seq3(ExactStr("a"), ExactStr("b"), ExactStr("e")),
				func(node *Seq3Node[string, string, string], _ Origin) string {
					return "e"
				})),
			Look(Map(Seq3(ExactStr("a"), ExactStr("b"), ExactStr("f")),
				func(node *Seq3Node[string, string, string], _ Origin) string {
					return "f"
				}))))

	expectedResults := []string{"c", "d", "e", "c", "f", "d"}
	for _, expected := range expectedResults {
		node, err := parser(ctx)
		assert.Equal(t, InvalidLook, ctx.lookVal)
		assert.NoError(t, err)
		assert.Equal(t, expected, node)
	}

	_, err := parser(ctx)
	assert.Error(t, err)
}
