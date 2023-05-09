package apcgen

import (
	"reflect"

	"github.com/tpillow/apc/pkg/apc"
)

var (
	runeParserCache  parserCache[rune]      = make(parserCache[rune])
	tokenParserCache parserCache[apc.Token] = make(parserCache[apc.Token])
)

func BuildParser[RT any](buildOpts *BuildOptions[rune]) apc.Parser[rune, RT] {
	rtType := reflect.TypeOf(*new(RT))
	buildCtx := newBuildContext(runeParserCache, buildOpts.ProvidedParsers)
	baseParser := buildRuneParserForType(buildCtx, rtType)
	parser := apc.CastTo[rune, any, RT](baseParser)
	parser = wrapWithSkipParsers(parser, buildOpts.SkipParsers)
	return parser
}

func BuildTokenizedParser[RT any](buildOpts *BuildOptions[apc.Token]) apc.Parser[apc.Token, RT] {
	rtType := reflect.TypeOf(*new(RT))
	buildCtx := newBuildContext(tokenParserCache, buildOpts.ProvidedParsers)
	baseParser := buildTokenParserForType(buildCtx, rtType)
	parser := apc.CastTo[apc.Token, any, RT](baseParser)
	parser = wrapWithSkipParsers(parser, buildOpts.SkipParsers)
	return parser
}
