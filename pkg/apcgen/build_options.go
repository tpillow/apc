package apcgen

import (
	"fmt"
	"reflect"

	"github.com/tpillow/apc/pkg/apc"
)

type BuildOptionFunc[CT any] func(opts *BuildOptions[CT])

type BuildOptions[CT any] struct {
	ProvidedParsers map[string]apc.Parser[CT, any]
	SkipParsers     []apc.Parser[CT, any]
}

func WithDefaultBuildOptions[CT any](buildFuncs ...BuildOptionFunc[CT]) *BuildOptions[CT] {
	opts := &BuildOptions[CT]{
		ProvidedParsers: make(map[string]apc.Parser[CT, any]),
		SkipParsers:     make([]apc.Parser[CT, any], 0),
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

func WithSkipParserOption[CT any](parser apc.Parser[CT, any]) BuildOptionFunc[CT] {
	return func(opts *BuildOptions[CT]) {
		opts.SkipParsers = append(opts.SkipParsers, parser)
	}
}

func WithBuildParserOption[RT any]() BuildOptionFunc[rune] {
	return func(opts *BuildOptions[rune]) {
		// TODO: make common function
		typeRef := reflect.TypeOf(*new(RT))
		typeName := typeRef.Name()
		if typeRef.Kind() == reflect.Pointer {
			typeName = typeRef.Elem().Name()
		}

		var parser apc.Parser[rune, RT]
		parserRef := apc.Ref(&parser)
		WithParserOption(typeName, apc.CastToAny(parserRef))(opts)
		parser = BuildParser[RT](opts)
	}
}

func WithBuildTokenizedParserOption[RT any]() BuildOptionFunc[apc.Token] {
	return func(opts *BuildOptions[apc.Token]) {
		// TODO: make common function
		typeRef := reflect.TypeOf(*new(RT))
		typeName := typeRef.Name()
		if typeRef.Kind() == reflect.Pointer {
			typeName = typeRef.Elem().Name()
		}

		var parser apc.Parser[apc.Token, RT]
		parserRef := apc.Ref(&parser)
		WithParserOption(typeName, apc.CastToAny(parserRef))(opts)
		parser = BuildTokenizedParser[RT](opts)
	}
}
