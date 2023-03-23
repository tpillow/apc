package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseErrorNoWrapsTypesIs(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "")
	pe := ParseErrExpectedButGot(ctx, false, true, nil)
	pec := ParseErrConsumedExpectedButGot(ctx, false, true, nil)

	assert.ErrorIs(t, pe, ErrParseErr)
	assert.ErrorIs(t, pec, ErrParseErrConsumed)

	assert.NotErrorIs(t, pe, ErrParseErrConsumed)
	assert.NotErrorIs(t, pec, ErrParseErr)

	assert.False(t, IsMustReturnParseErr(pe))
	assert.True(t, IsMustReturnParseErr(pec))
}

func TestParseErrorConsumedWrapsTypesIs(t *testing.T) {
	ctx := NewStringContext(testStringOrigin, "")
	pe := ParseErrExpectedButGot(ctx, false, true, nil)
	pec := ParseErrConsumedExpectedButGot(ctx, false, true, pe)

	assert.ErrorIs(t, pe, ErrParseErr)
	assert.ErrorIs(t, pec, ErrParseErrConsumed)

	assert.NotErrorIs(t, pe, ErrParseErrConsumed)
	assert.ErrorIs(t, pec, ErrParseErr)

	assert.False(t, IsMustReturnParseErr(pe))
	assert.True(t, IsMustReturnParseErr(pec))
}
