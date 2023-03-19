package apc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapParserMatches(t *testing.T) {
	RunBasicParserMatchTest(t, Map(Exact("hi"), func(node Node) Node {
		return 555
	}), "hi#hi", 555, 555)
}

func TestBind(t *testing.T) {
	ctx := NewStringContext(strOriginName, "hi")

	node, err := Bind(Exact("hi"), 555).Parse(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 555, node)
}
