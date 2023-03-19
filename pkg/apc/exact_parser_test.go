package apc

import (
	"testing"
)

func TestExactParserMatches(t *testing.T) {
	RunBasicParserMatchTest(t, Exact("hello!"), "hello!$hello!", "hello!", "hello!")
}
