package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOneParserMatches(t *testing.T) {
	p := One(Exact("hi"), Exact("bye"))
	RunBasicParserMatchTest(t, p,
		"hi#bye", "hi", "bye")
	ctx := NewStringContext(strOriginName, "#")
	_, err := p.Parse(ctx)
	assert.Error(t, err)
}
