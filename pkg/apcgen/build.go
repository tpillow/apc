package apcgen

import (
	"github.com/tpillow/apc/pkg/apc"
)

var (
	runeParserCache  parserCache[rune]      = make(parserCache[rune])
	tokenParserCache parserCache[apc.Token] = make(parserCache[apc.Token])
)

func BuildParser[RT any](buildOpts *BuildOptions[rune]) apc.Parser[rune, RT] {
	resultType := reflectTypeOf[RT]()
	buildCtx := newBuildContext(runeParserCache, buildOpts.ProvidedParsers)
	baseParser := buildRuneParserForType(buildCtx, resultType)
	parser := apc.CastTo[rune, any, RT](baseParser)
	parser = wrapWithSkipParsers(parser, buildOpts.SkipParsers)
	return parser
}

func BuildTokenizedParser[RT any](buildOpts *BuildOptions[apc.Token]) apc.Parser[apc.Token, RT] {
	resultType := reflectTypeOf[RT]()
	buildCtx := newBuildContext(tokenParserCache, buildOpts.ProvidedParsers)
	baseParser := buildTokenParserForType(buildCtx, resultType)
	parser := apc.CastTo[apc.Token, any, RT](baseParser)
	parser = wrapWithSkipParsers(parser, buildOpts.SkipParsers)
	return parser
}
