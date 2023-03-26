package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFuncOptionMustParseToEOF(t *testing.T) {
	ctx := NewRuneContextFromStr(testStringOrigin, "hihiho")
	p := ExactStr("hi")

	node, err := Parse[rune](ctx, p, ParseConfig{MustParseToEOF: false})
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	_, err = Parse[rune](ctx, p, ParseConfig{MustParseToEOF: true})
	assert.ErrorIs(t, err, ErrParseErr)
}
