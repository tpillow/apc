package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyContext(t *testing.T) {
	ctx := NewStringContext("<string>", "")
	assert.True(t, ctx.IsEof())
	assert.Panics(t, func() { ctx.Peek() })
	assert.Panics(t, func() { ctx.Pop() })
}

func TestNonEmptyContext(t *testing.T) {
	ctx := NewStringContext("<string>", "ab")
	assert.False(t, ctx.IsEof())
	assert.Equal(t, 'a', ctx.Peek())
	assert.Equal(t, 'a', ctx.Pop())
	assert.False(t, ctx.IsEof())
	assert.Equal(t, 'b', ctx.Peek())
	assert.Equal(t, 'b', ctx.Pop())
	assert.True(t, ctx.IsEof())
	assert.Panics(t, func() { ctx.Peek() })
	assert.Panics(t, func() { ctx.Pop() })
}
