package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/apc/pkg/apc"
)

func TestRuneParserCaptureString(t *testing.T) {
	type Person struct {
		Name string `apc:"'name' $'Tommy'"`
		Age  string `apc:"'age' $'29'"`
	}

	parser, err := BuildRuneParser[Person](DefaultBuildOptions)
	assert.NoError(t, err)

	ctx := apc.NewStringContext(testOriginName, `name Tommy age 29`)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Person{"Tommy", "29"}, node)
}
