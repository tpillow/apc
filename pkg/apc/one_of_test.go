package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOneOfParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "###hi##")
	p1 := Regex("", "#+")
	p2 := ExactStr("hi")
	p := OneOf("", p1, p2)

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "###", node)

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "##", node)

	_, err = ctx.Peek(0, 1)
	assert.ErrorIs(t, err, ErrEOF)
}
