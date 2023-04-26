package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeqParser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "55,66#")
	p := Seq(CastToAny(IntParser), CastToAny(ExactStr(",")))

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []any{int64(55), ","}, node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErrConsumed)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)
	ctx.Consume(1)

	_, err = ctx.Peek(0, 1)
	assert.ErrorIs(t, err, ErrEOF)
}

func TestSeq2Parser(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "55,66#")
	p := Seq2(IntParser, ExactStr(","))

	node, err := p(ctx)
	assert.NoError(t, err)
	assert.Equal(t, &Seq2Node[int64, string]{
		Result1: 55,
		Result2: ",",
	}, node)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErrConsumed)

	_, err = p(ctx)
	assert.ErrorIs(t, err, ErrParseErr)
	ctx.Consume(1)

	_, err = ctx.Peek(0, 1)
	assert.ErrorIs(t, err, ErrEOF)
}
