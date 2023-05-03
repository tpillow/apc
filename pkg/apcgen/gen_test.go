package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/apc/pkg/apc"
)

func TestRuneParserCaptureString(t *testing.T) {
	type Person struct {
		Name string `apc:"'person' $'Tommy'"`
		Age  string `apc:"$'29'"`
	}

	parser, err := BuildRuneParser[Person](DefaultBuildOptions)
	assert.NoError(t, err)

	ctx := apc.NewStringContext(testOriginName, `person Tommy 29`)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Person{"Tommy", "29"}, node)
}
