package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	ctx := NewStringContext("<string>", "ab")
	result, err := NewBuilder(Exact('a')).Skip(Exact('b')).Build().ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 'a', result)
}

func TestThen(t *testing.T) {
	ctx := NewStringContext("<string>", "ab")
	result, err := NewBuilder(Exact('a')).Then(Exact('b')).Build().ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 'b', result)
}
