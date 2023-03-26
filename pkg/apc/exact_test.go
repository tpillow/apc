package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExactParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "hiho")
	p := ExactStr("hi")

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, []rune{'h'}, r)
}
