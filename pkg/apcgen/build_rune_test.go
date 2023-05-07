package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/apc/pkg/apc"
)

func TestParserCaptureStringAndRegex(t *testing.T) {
	type Person struct {
		Name string `apc:"'person' $regex('[a-zA-Z]+')"`
		Age  string `apc:"$'29'"`
	}

	parser := BuildParser[Person](WithDefaultBuildOptions[rune](), true)

	ctx := apc.NewStringContext(testOriginName, `person Tommy 29`)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Person{"Tommy", "29"}, node)
}

func TestSliceCaptureString(t *testing.T) {
	type Obj struct {
		Values []string `apc:"$'ha'*"`
	}

	parser := BuildParser[Obj](WithDefaultBuildOptions[rune](), true)

	ctx := apc.NewStringContext(testOriginName, `ha ha ha ha`)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Obj{
		Values: []string{"ha", "ha", "ha", "ha"},
	}, node)
}

func TestSliceCaptureStruct(t *testing.T) {
	type NameObj struct {
		Name string `apc:"$regex('[a-zA-Z0-9]+')"`
	}
	type Obj struct {
		NameObjs []*NameObj `apc:"$.*"`
	}

	parser := BuildParser[Obj](WithDefaultBuildOptions[rune](), true)

	ctx := apc.NewStringContext(testOriginName, `Name1 Name2 Name3`)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Obj{
		NameObjs: []*NameObj{{Name: "Name1"}, {Name: "Name2"}, {Name: "Name3"}},
	}, node)
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

		BoolFromStrT   bool `apc:"$'true'"`
		BoolFromStrF   bool `apc:"$'false'"`
		BoolFromMaybeT bool `apc:"$'maybe'?"`
		BoolFromMaybeF bool `apc:"$'nonexistent'?"`
	}
	parser := BuildParser[Obj](WithDefaultBuildOptions[rune](), true)

	ctx := apc.NewStringContext(testOriginName, `strVal 51 52 53 54 55 56 57 58 59 60 71.1 72.2 true false maybe`)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Obj{
		StrVal:         "strVal",
		IntVal:         51,
		Int8Val:        52,
		Int16Val:       53,
		Int32Val:       54,
		Int64Val:       55,
		UIntVal:        56,
		UInt8Val:       57,
		UInt16Val:      58,
		UInt32Val:      59,
		UInt64Val:      60,
		Float32Val:     71.1,
		Float64Val:     72.2,
		BoolFromStrT:   true,
		BoolFromStrF:   false,
		BoolFromMaybeT: true,
		BoolFromMaybeF: false,
	}, node)
}
