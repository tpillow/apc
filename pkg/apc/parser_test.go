package apc

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExactParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "hiho")
	p := Exact("hi")

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, "h", r)
}

func TestRegexParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "###_##")
	p := Regex("", "#+")

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "###", node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, "_", r)
}

func TestAnyParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "###hi##")
	p1 := Regex("", "#+")
	p2 := Exact("hi")
	p := OneOf("", p1, p2)

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "###", node)

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "##", node)

	_, err = ctx.Peek(0, 1)
	assert.ErrorIs(t, err, ErrEOF)
}

func TestMapParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "342_")
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
	assert.Equal(t, "_", r)
}

func TestSkipParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, " \t\nhi\n\n\t  hi_")
	p := Skip(MapToAny(WhitespaceParser), Exact("hi"))

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, "_", r)
}

func TestUnskipParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, " hi hi_")
	wsp := MapToAny(WhitespaceParser)
	ctx.AddSkipParser(wsp)
	p := Unskip(wsp, Exact("hi"))

	_, err := p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)
	ctx.Consume(1)

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)
	ctx.Consume(1)

	node, err = p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	r, err := ctx.Peek(0, 1)
	assert.NoError(t, err)
	assert.Equal(t, "_", r)
}
