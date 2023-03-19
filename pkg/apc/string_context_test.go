package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const strOriginName = "<str>"

func TestStringContextOriginStartsCorrect(t *testing.T) {
	ctx := NewStringContext(strOriginName, "abc")
	assert.Equal(t, ctx.GetOrigin(), Origin{
		Name:    strOriginName,
		LineNum: 1,
		ColNum:  1,
	})
}

func TestStringContextEmpty(t *testing.T) {
	ctx := NewStringContext(strOriginName, "")
	_, err := ctx.PeekRune(0)
	assert.ErrorIs(t, err, &EOFError{})
	_, err = ctx.PeekRune(1)
	assert.ErrorIs(t, err, &EOFError{})
	_, err = ctx.ConsumeRune()
	assert.ErrorIs(t, err, &EOFError{})
	assert.Equal(t, ctx.GetOrigin(), Origin{
		Name:    strOriginName,
		LineNum: 1,
		ColNum:  1,
	})
}

func TestStringContextLineCountConsumeToEnd(t *testing.T) {
	ctx := NewStringContext(strOriginName, "ab\nc")

	r, err := ctx.PeekRune(0)
	assert.NoError(t, err)
	assert.Equal(t, 'a', r)
	assert.Equal(t, ctx.GetOrigin(), Origin{
		Name:    strOriginName,
		LineNum: 1,
		ColNum:  1,
	})

	r, err = ctx.PeekRune(3)
	assert.NoError(t, err)
	assert.Equal(t, 'c', r)
	assert.Equal(t, ctx.GetOrigin(), Origin{
		Name:    strOriginName,
		LineNum: 1,
		ColNum:  1,
	})

	_, err = ctx.PeekRune(4)
	assert.ErrorIs(t, err, &EOFError{})

	r, err = ctx.ConsumeRune()
	assert.NoError(t, err)
	assert.Equal(t, 'a', r)
	assert.Equal(t, ctx.GetOrigin(), Origin{
		Name:    strOriginName,
		LineNum: 1,
		ColNum:  2,
	})

	r, err = ctx.ConsumeRune()
	assert.NoError(t, err)
	assert.Equal(t, 'b', r)
	assert.Equal(t, ctx.GetOrigin(), Origin{
		Name:    strOriginName,
		LineNum: 1,
		ColNum:  3,
	})

	r, err = ctx.ConsumeRune()
	assert.NoError(t, err)
	assert.Equal(t, '\n', r)
	assert.Equal(t, ctx.GetOrigin(), Origin{
		Name:    strOriginName,
		LineNum: 2,
		ColNum:  1,
	})

	r, err = ctx.ConsumeRune()
	assert.NoError(t, err)
	assert.Equal(t, 'c', r)
	assert.Equal(t, ctx.GetOrigin(), Origin{
		Name:    strOriginName,
		LineNum: 2,
		ColNum:  2,
	})

	_, err = ctx.ConsumeRune()
	assert.ErrorIs(t, err, &EOFError{})
}
