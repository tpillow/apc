package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	ctx := NewSliceContext([]rune("test"))
	result, err := Match('t').Parse(ctx)
	assert.Equal(t, result, 't')
	assert.Nil(t, err)
	_, err = Match('t').Parse(ctx)
	assert.Error(t, err)
}
