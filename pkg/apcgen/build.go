package apcgen

import (
	"reflect"

	"github.com/tpillow/apc/pkg/apc"
)

func BuildParser[RT any](buildOpts RuneBuildOptions, providedParsers map[string]apc.Parser[rune, any]) apc.Parser[rune, *RT] {
	rtType := reflect.TypeOf(new(RT))

	buildCtx := newBuildContext(providedParsers)
	baseParser := buildRuneParserForType(buildCtx, rtType)
	parser := apc.CastTo[rune, any, *RT](baseParser)

	if buildOpts.SkipWhitespace {
		parser = apc.Skip(
			apc.CastToAny(apc.WhitespaceParser),
			parser,
		)
	}

	return parser
}

func BuildTokenizedParser[RT any](providedParsers map[string]apc.Parser[apc.Token, any]) apc.Parser[apc.Token, *RT] {
	rtType := reflect.TypeOf(new(RT))

	buildCtx := newBuildContext(providedParsers)
	baseParser := buildTokenParserForType(buildCtx, rtType)
	return apc.CastTo[apc.Token, any, *RT](baseParser)
}
