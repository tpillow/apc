package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testStringOrigin = "<origin>"

func TestStringContextEmptyInput(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, []rune(""))
	assert.Equal(t, Origin{
		Name:    testStringOrigin,
		LineNum: 1,
		ColNum:  1,
	}, ctx.GetCurOrigin())

	assert := func(val string, err error) {
		assert.Equal(t, "", val)
		assert.ErrorIs(t, err, ErrEOF)
	}

	val, err := ctx.Peek(0, 1)
	assert(val, err)
	val, err = ctx.Peek(2, 10)
	assert(val, err)

	val, err = ctx.Consume(1)
	assert(val, err)
	val, err = ctx.Consume(10)
	assert(val, err)
}

func TestStringContextBasic(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, []rune("ab\ncd"))

	assert := func(val string, exp string, line int, col int, err error, expErr bool) {
		if expErr {
			assert.ErrorIs(t, err, ErrEOF)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, exp, val)
		assert.Equal(t, Origin{
			Name:    testStringOrigin,
			LineNum: line,
			ColNum:  col,
		}, ctx.GetCurOrigin())
	}

	val, err := ctx.Peek(2, 3)
	assert(val, "\ncd", 1, 1, err, false)
	val, err = ctx.Peek(0, 7)
	assert(val, "ab\ncd", 1, 1, err, true)

	val, err = ctx.Consume(1)
	assert(val, "a", 1, 2, err, false)
	val, err = ctx.Consume(2)
	assert(val, "b\n", 2, 1, err, false)
	val, err = ctx.Consume(3)
	assert(val, "cd", 2, 3, err, true)
}
