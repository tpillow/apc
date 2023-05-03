package apcgen

import (
	"github.com/stretchr/testify/assert"
	"github.com/tpillow/apc/pkg/apc"
	"testing"
)

func TestParserCaptureString(t *testing.T) {
	type Person struct {
		Name string `apc:"'person' $'Tommy'"`
		Age  string `apc:"$'29'"`
	}

	parser := BuildParser[*Person](DefaultBuildOptions, map[string]apc.Parser[rune, any]{})

	ctx := apc.NewStringContext(testOriginName, `person Tommy 29`)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Person{"Tommy", "29"}, node)
}

func TestParserIntrinsicConversions(t *testing.T) {
	type Obj struct {
		StrVal string `apc:"$'strVal'"`

		IntVal   int   `apc:"$'51'"`
		Int8Val  int8  `apc:"$'52'"`
		Int16Val int16 `apc:"$'53'"`
		Int32Val int32 `apc:"$'54'"`
		Int64Val int64 `apc:"$'55'"`

		UIntVal   uint   `apc:"$'56'"`
		UInt8Val  uint8  `apc:"$'57'"`
		UInt16Val uint16 `apc:"$'58'"`
		UInt32Val uint32 `apc:"$'59'"`
		UInt64Val uint64 `apc:"$'60'"`

		Float32Val float32 `apc:"$'71.1'"`
		Float64Val float64 `apc:"$'72.2'"`

		BoolValT  bool `apc:"$'true'"`
		BoolValF  bool `apc:"$'false'"`
		BoolValIT bool `apc:"$'1'"`
		BoolValIF bool `apc:"$'0'"`
	}
	parser := BuildParser[*Obj](DefaultBuildOptions, map[string]apc.Parser[rune, any]{})

	ctx := apc.NewStringContext(testOriginName, `strVal 51 52 53 54 55 56 57 58 59 60 71.1 72.2 true false 1 0`)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Obj{
		StrVal:     "strVal",
		IntVal:     51,
		Int8Val:    52,
		Int16Val:   53,
		Int32Val:   54,
		Int64Val:   55,
		UIntVal:    56,
		UInt8Val:   57,
		UInt16Val:  58,
		UInt32Val:  59,
		UInt64Val:  60,
		Float32Val: 71.1,
		Float64Val: 72.2,
		BoolValT:   true,
		BoolValF:   false,
		BoolValIT:  true,
		BoolValIF:  false,
	}, node)
}
