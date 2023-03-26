package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFuncOptionMustParseToEOF(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "hihiho")
	p := ExactStr("hi")

	node, err := Parse[string](ctx, p, ParseConfig{MustParseToEOF: false})
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	_, err = Parse[string](ctx, p, ParseConfig{MustParseToEOF: true})
	assert.ErrorIs(t, err, ErrParseErr)
}
