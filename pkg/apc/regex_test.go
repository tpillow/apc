package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "###_##")
	p := Regex("", "#+")

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "###", node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, "_", r)
}
