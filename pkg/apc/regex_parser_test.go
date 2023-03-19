package apc

import "testing"

func TestRegexParserMatches(t *testing.T) {
	RunBasicParserMatchTest(t, Regex("\\d+"), "243#5", "243", "5")
}
