package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZeroOrOneParserMatches(t *testing.T) {
	ctx := NewStringContext(strOriginName, "hi#")
	p := ZeroOrOne(Exact("hi"))

	node, err := p.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "hi", node)

	node, err = p.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, nil, node)
}

func TestZeroOrMoreParserMatches(t *testing.T) {
	ctx := NewStringContext(strOriginName, "hihi#")
	p := ZeroOrMore(Exact("hi"))

	node, err := p.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []Node{"hi", "hi"}, node)

	node, err = p.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []Node{}, node)
}

func TestOneOrMoreParserMatches(t *testing.T) {
	ctx := NewStringContext(strOriginName, "hihi#")
	p := OneOrMore(Exact("hi"))

	node, err := p.Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, []Node{"hi", "hi"}, node)

	_, err = p.Parse(ctx)
	assert.Error(t, err)
}
