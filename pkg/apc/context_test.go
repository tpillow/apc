package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeekNRunes(t *testing.T) {
	ctx := NewStringContext(strOriginName, "abc")

	val, err := PeekNRunes(ctx, 1, 2)
	assert.NoError(t, err)
	assert.Equal(t, "bc", val)

	val, err = PeekNRunes(ctx, 0, 4)
	assert.ErrorIs(t, err, &EOFError{})
	assert.Equal(t, "abc", val)
}

func TestConsumeNRunes(t *testing.T) {
	ctx := NewStringContext(strOriginName, "abcd")

	val, err := ConsumeNRunes(ctx, 2)
	assert.NoError(t, err)
	assert.Equal(t, "ab", val)

	val, err = ConsumeNRunes(ctx, 3)
	assert.ErrorIs(t, err, &EOFError{})
	assert.Equal(t, "cd", val)
}
