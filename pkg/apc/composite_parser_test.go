package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	ctx := NewStringContext("ab")
	result, err := Exact('a').Skip(Exact('b')).ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 'a', result)
}

func TestThen(t *testing.T) {
	ctx := NewStringContext("ab")
	result, err := Exact('a').Then(Exact('b')).ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 'b', result)
}
