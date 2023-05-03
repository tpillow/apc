package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRangeParser(t *testing.T) {
	p := Range(2, 3, ExactStr("hi"))

	ctx := NewStringContext(testStringOrigin, "hi")
	_, err := p(ctx)
	assert.ErrorIs(t, err, ErrParseErrConsumed)

	ctx = NewStringContext(testStringOrigin, "hihihihi")
	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"hi", "hi", "hi"}, node)
	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, []rune{'h'}, r)

	ctx = NewStringContext(testStringOrigin, "hihi")
	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"hi", "hi"}, node)

	ctx = NewStringContext(testStringOrigin, "hihihi")
	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []string{"hi", "hi", "hi"}, node)

	ctx = NewStringContext(testStringOrigin, "__")
	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)
}

func TestMaybeParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "hibye")
	p := Maybe(ExactStr("hi"))

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.False(t, node.IsNil())
	assert.Equal(t, "hi", node.Value())

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.True(t, node.IsNil())

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, []rune{'b'}, r)
}

func TestOneOrMoreParserWithSeq2(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "#$#$")
	p := OneOrMore(Seq(ExactStr("#"), ExactStr("$")))

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{
		{"#", "$"},
		{"#", "$"},
	}, node)

	ctx = NewStringContext(testStringOrigin, "#$#$#")
	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErrConsumed)
}

func TestOneOrMoreSeparatedParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "55,66")
	p := OneOrMoreSeparated(IntParser, ExactStr(","))

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []int64{55, 66}, node)

	_, err = ctx.Peek(0, 1)
	assert.ErrorIs(t, err, ErrEOF)

	ctx = NewStringContext(testStringOrigin, "55,66,")
	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErrConsumed)
}
