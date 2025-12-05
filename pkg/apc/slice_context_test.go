package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyContext(t *testing.T) {
	ctx := NewStringContext("")
	assert.True(t, ContextIsEof(ctx))
	assert.True(t, IsEofValue(ctx.Peek()))
	assert.True(t, IsEofValue(ctx.Pop()))
}

func TestNonEmptyContext(t *testing.T) {
	ctx := NewStringContext("ab")
	assert.False(t, ContextIsEof(ctx))
	assert.Equal(t, 'a', ctx.Peek())
	assert.Equal(t, 'a', ctx.Pop())
	assert.False(t, ContextIsEof(ctx))
	assert.Equal(t, 'b', ctx.Peek())
	assert.Equal(t, 'b', ctx.Pop())
	assert.True(t, ContextIsEof(ctx))
	assert.True(t, IsEofValue(ctx.Peek()))
	assert.True(t, IsEofValue(ctx.Pop()))
}
