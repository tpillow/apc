package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "#hello##hello")
	var value Parser[rune, any]
	var valueRef = Ref(&value)
	var hashValue = Seq("", MapToAny(ExactStr("#")), valueRef)
	value = OneOf("", MapToAny(ExactStr("hello")), MapToAny(hashValue))

	node, err := valueRef(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []any{"#", "hello"}, node)

	node, err = hashValue(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []any{"#", []any{"#", "hello"}}, node)

	_, err = valueRef(ctx)
	assert.ErrorIs(t, err, ErrParseErr)
}
