package apcgen

import (
	"reflect"

	"github.com/tpillow/apc/pkg/apc"
)

func BuildParser[RT any](buildOpts *BuildOptions[rune], skipWhitespace bool) apc.Parser[rune, *RT] {
	rtType := reflect.TypeOf(new(RT))

	buildCtx := newBuildContext(buildOpts.ProvidedParsers)
	baseParser := buildRuneParserForType(buildCtx, rtType)
	parser := apc.CastTo[rune, any, *RT](baseParser)

	if skipWhitespace {
		parser = apc.Skip(
			apc.CastToAny(apc.WhitespaceParser),
			parser,
		)
	}

	return parser
}

func BuildTokenizedParser[RT any](buildOpts *BuildOptions[apc.Token]) apc.Parser[apc.Token, *RT] {
	rtType := reflect.TypeOf(new(RT))

	buildCtx := newBuildContext(buildOpts.ProvidedParsers)
	baseParser := buildTokenParserForType(buildCtx, rtType)
	return apc.CastTo[apc.Token, any, *RT](baseParser)
}
