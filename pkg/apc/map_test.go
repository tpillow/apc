package apc

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapParser(t *testing.T) {
	ctx := NewRuneContextFromStr(testStringOrigin, "342_")
	p := Map(Regex("", "\\d+"), func(node string) int64 {
		val, err := strconv.ParseInt(node, 10, 64)
		assert.NoError(t, err)
		return val
	})

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(342), node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, []rune{'_'}, r)
}

func TestBindParser(t *testing.T) {
	ctx := NewRuneContextFromStr(testStringOrigin, "hi_")
	p := Bind(ExactStr("hi"), 55)

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 55, node)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, []rune{'_'}, r)
}
