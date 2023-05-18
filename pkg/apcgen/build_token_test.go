package apcgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tpillow/apc/pkg/apc"
)

func TestTokenParserBasic1(t *testing.T) {
	type AssignExpr struct {
		Name string `apc:"look($'Ident' token('Ident', '='))"`
		Expr string `apc:"$'Ident' token('Ident', ';')"`
	}
	type TypeExpr struct {
		Name string `apc:"look($'Ident' token('Ident', 'type'))"`
		_    string `apc:"token('Ident', ';')"`
	}
	type Result struct {
		Things []any `apc:"$(AssignExpr | TypeExpr)*"`
	}

	var TokenTypeIdent apc.TokenType = "Ident"
	lexParser := apc.Skip(
		apc.CastToAny(apc.WhitespaceParser),
		apc.Any(
			apc.BindToToken(apc.Regex("[^\\s]+"), TokenTypeIdent),
		),
	)

	parser := BuildTokenizedParser[Result](WithDefaultBuildOptions(
		WithBuildTokenizedParserOption[AssignExpr](),
		WithBuildTokenizedParserOption[TypeExpr](),
	))

	ctx := apc.NewStringContext(testOriginName, `name = value ; some type ; another type ; name2 = hahaha ;`)
	lexer := apc.NewParseReader[rune](ctx, lexParser)
	lexerCtx := apc.NewReaderContext[apc.Token](lexer)

	node, err := apc.Parse[apc.Token](lexerCtx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, Result{
		Things: []any{
			AssignExpr{Name: "name", Expr: "value"},
			TypeExpr{Name: "some"},
			TypeExpr{Name: "another"},
			AssignExpr{Name: "name2", Expr: "hahaha"},
		},
	}, node)
}

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

	parser := BuildTokenizedParser[*Person](WithDefaultBuildOptions[apc.Token]())

	ctx := apc.NewStringContext(testOriginName, `person Tommy 29`)
	lexer := apc.NewParseReader[rune](ctx, lexParser)
	lexerCtx := apc.NewReaderContext[apc.Token](lexer)

	node, err := apc.Parse[apc.Token](lexerCtx, parser, apc.DefaultParseConfig)
	assert.NoError(t, err)
	assert.Equal(t, &Person{"Tommy", 29}, node)
}
