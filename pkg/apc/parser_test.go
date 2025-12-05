package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	ctx := NewStringContext("4Hfour")
	result, err := Int64.Skip(Exact('H')).Generate(func(rawValue any) *Parser {
		value := rawValue.(int64)
		return AnyChar.Times(int(value), int(value))
	}).ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []any{"f", "o", "u", "r"}, result)
}
