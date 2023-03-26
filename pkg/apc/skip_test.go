package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkipParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, " \t\nhi\n\n\t  hi_")
	p := Skip(MapToAny(WhitespaceParser), ExactStr("hi"))

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, "_", r)
}

func TestUnskipParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, " hi hi_")
	wsp := MapToAny(WhitespaceParser)
	ctx.AddSkipParser(wsp)
	p := Unskip(wsp, ExactStr("hi"))

	_, err := p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)
	ctx.Consume(1)

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)
	ctx.Consume(1)

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, "_", r)
}
