package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExactStr(t *testing.T) {
	ctx := NewStringContext("abc")
	result, err := ExactStr("abc").ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "abc", result)
}
