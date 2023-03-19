package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Input string should be in exact format of `<expectA><nonMatchRune><expectB><EOF>`
func RunBasicParserMatchTest(t *testing.T, parser Parser, data string, expectA Node, expectB Node) {
	ctx := NewStringContext(strOriginName, data)

	node, err := parser.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectA, node)

	_, err = parser.Parse(ctx)
	assert.Error(t, err)
	ctx.ConsumeRune()

	node, err = parser.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectB, node)

	_, err = ctx.PeekRune(0)
	assert.ErrorIs(t, err, &EOFError{})
}
