package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloat(t *testing.T) {
	ctx := NewStringContext("034.022")
	result, err := Float.ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 34.022, result)
}
