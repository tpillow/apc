package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	ctx := NewStringContext("4Hfour")
	result, err := NewBuilder(Int64).Skip(Exact('H')).Generate(func(rawValue any) *Parser {
		value := rawValue.(int64)
		return NewBuilder(AnyChar).Times(int(value), int(value)).Build()
	}).Build().ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []any{"f", "o", "u", "r"}, result)
}
