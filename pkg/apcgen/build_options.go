package apcgen

import (
	"fmt"
	"reflect"

	"github.com/tpillow/apc/pkg/apc"
)

type BuildOptionFunc[CT any] func(opts *BuildOptions[CT])

type BuildOptions[CT any] struct {
	ProvidedParsers map[string]apc.Parser[CT, any]
}

func WithDefaultBuildOptions[CT any](buildFuncs ...BuildOptionFunc[CT]) *BuildOptions[CT] {
	opts := &BuildOptions[CT]{
		ProvidedParsers: make(map[string]apc.Parser[CT, any]),
	}
	for _, buildFunc := range buildFuncs {
		buildFunc(opts)
	}
	return opts
}

func WithParserOption[CT any](name string, parser apc.Parser[CT, any]) BuildOptionFunc[CT] {
	return func(opts *BuildOptions[CT]) {
		if _, has := opts.ProvidedParsers[name]; has {
			panic(fmt.Sprintf("cannot use WithParserOption: name '%v' already specified", name))
		}
		opts.ProvidedParsers[name] = parser
	}
}

func WithBuildParserOption[RT any](skipWhitespace bool) BuildOptionFunc[rune] {
	return func(opts *BuildOptions[rune]) {
		typeName := reflect.TypeOf(new(RT)).Elem().Name()
		var parser *apc.Parser[rune, *RT]
		parserRef := apc.Ref(parser)
		WithParserOption(typeName, apc.CastToAny(parserRef))
		*parser = BuildParser[RT](opts, skipWhitespace)
	}
}

func WithBuildTokenizedParserOption[RT any]() BuildOptionFunc[apc.Token] {
	return func(opts *BuildOptions[apc.Token]) {
		typeName := reflect.TypeOf(new(RT)).Elem().Name()
		var parser *apc.Parser[apc.Token, *RT]
		parserRef := apc.Ref(parser)
		WithParserOption(typeName, apc.CastToAny(parserRef))
		*parser = BuildTokenizedParser[RT](opts)
	}
}
