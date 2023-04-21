package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLook(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "(hi)(bye)((hi))(((bye)))")

	var oneParser Parser[rune, string]
	parser := Map(
		Seq3("",
			ExactStr("("),
			Ref(&oneParser),
			ExactStr(")")),
		func(node *Seq3Node[string, string, string], origin Origin) string {
			return node.Result2
		})
	oneParser = OneOf("", ExactStr("hi"), ExactStr("bye"), parser)

	expectedResults := []string{"hi", "bye", "hi", "bye"}
	for _, expected := range expectedResults {
		node, err := parser(ctx)
		assert.Equal(t, 0, len(ctx.lookStack))
		assert.NoError(t, err)
		assert.Equal(t, expected, node)
	}

	_, err := parser(ctx)
	assert.Error(t, err)
}
