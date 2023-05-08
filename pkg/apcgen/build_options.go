package apcgen

import (
	"fmt"

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
