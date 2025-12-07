package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat(t *testing.T) {
	ctx := NewStringContext("<string>", "034.022")
	result, err := Float64.ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 34.022, result)
}
