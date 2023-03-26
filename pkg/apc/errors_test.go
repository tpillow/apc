package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseErrorNoWrapsTypesIs(t *testing.T) {
	ctx := NewRuneContextFromStr(testStringOrigin, "")
	pe := ParseErrExpectedButGot[rune](ctx, false, true, nil)
	pec := ParseErrConsumedExpectedButGot[rune](ctx, false, true, nil)

	assert.ErrorIs(t, pe, ErrParseErr)
	assert.ErrorIs(t, pec, ErrParseErrConsumed)

	assert.NotErrorIs(t, pe, ErrParseErrConsumed)
	assert.NotErrorIs(t, pec, ErrParseErr)

	assert.False(t, IsMustReturnParseErr(pe))
	assert.True(t, IsMustReturnParseErr(pec))
}

func TestParseErrorConsumedWrapsTypesIs(t *testing.T) {
	ctx := NewRuneContextFromStr(testStringOrigin, "")
	pe := ParseErrExpectedButGot[rune](ctx, false, true, nil)
	pec := ParseErrConsumedExpectedButGot[rune](ctx, false, true, pe)

	assert.ErrorIs(t, pe, ErrParseErr)
	assert.ErrorIs(t, pec, ErrParseErrConsumed)

	assert.NotErrorIs(t, pe, ErrParseErrConsumed)
	assert.ErrorIs(t, pec, ErrParseErr)

	assert.False(t, IsMustReturnParseErr(pe))
	assert.True(t, IsMustReturnParseErr(pec))
}
