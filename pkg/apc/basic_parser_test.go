package apc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEofPositive(t *testing.T) {
	ctx := NewStringContext("")
	result, err := Eof().ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, nil, result)
}

func TestEofNegative(t *testing.T) {
	ctx := NewStringContext("a")
	_, err := Eof().ParseToEof(ctx)
	assert.Error(t, err)
}

func TestExact(t *testing.T) {
	ctx := NewStringContext("a")
	result, err := Exact('a').ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 'a', result)
}

func TestSucceed(t *testing.T) {
	ctx := NewStringContext("")
	result, err := Succeed(6).ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 6, result)
}

func TestFail(t *testing.T) {
	ctx := NewStringContext("")
	expectError := fmt.Errorf("error")
	_, err := Fail(expectError).ParseToEof(ctx)
	assert.Error(t, err)
}

func TestSeq(t *testing.T) {
	ctx := NewStringContext("abc")
	result, err := Seq(Exact('a'), Exact('b'), Exact('c')).ParseToEof(ctx)
	assert.NoError(t, err)
	results, ok := result.([]any)
	assert.True(t, ok)
	assert.Len(t, results, 3)
	assert.Equal(t, 'a', results[0])
	assert.Equal(t, 'b', results[1])
	assert.Equal(t, 'c', results[2])
}

func TestAnyOf(t *testing.T) {
	ctx := NewStringContext("a")
	result, err := AnyOf(Exact('b'), Exact('a'), Exact('c')).ParseToEof(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 'a', result)
}
