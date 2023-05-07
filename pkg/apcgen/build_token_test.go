package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/apc/pkg/apc"
)

func TestTokenParserCaptureBasic(t *testing.T) {
	type Person struct {
		Name string `apc:"'Ident:person' $'Ident'"`
		Age  int    `apc:"$'Int'"`
	}

	var TokenTypeIdent apc.TokenType = "Ident"
	var TokenTypeInt apc.TokenType = "Int"

	lexParser := apc.Skip(
		apc.CastToAny(apc.WhitespaceParser),
		apc.Any(
			apc.BindToToken(apc.IdentifierParser, TokenTypeIdent),
			apc.BindToToken(apc.IntParser, TokenTypeInt),
		),
	)

	parser := BuildTokenizedParser[Person](WithDefaultBuildOptions[apc.Token]())

	ctx := apc.NewStringContext(testOriginName, `person Tommy 29`)
	lexer := apc.NewParseReader[rune](ctx, lexParser)
	lexerCtx := apc.NewReaderContext[apc.Token](lexer)

	node, err := apc.Parse[apc.Token](lexerCtx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Person{"Tommy", 29}, node)
}
