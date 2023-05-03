package main

import (
	"fmt"

	"github.com/tpillow/apc/pkg/apc"
	"github.com/tpillow/apc/pkg/apcgen"
)

type Dir struct {
	Entries []*DirEntry `apc:"'Dir' '{' $(.*) '}'"`
}

type DirEntry struct {
	Name   string `apc:"'Entry' '{' $StrParser"`
	Id     int    `apc:"$regex('[0-9]+')"`
	SubDir *Dir   `apc:"$.? '}'"`
}

func main() {
	parser := apcgen.BuildParser[*Dir](
		apcgen.DefaultBuildOptions,
		map[string]apc.Parser[rune, any]{
			"StrParser": apc.CastToAny(apc.DoubleQuotedStringParser),
		},
	)

	input := `Dir { Entry { "Name1" 1 } Entry { "Name2" 2 Dir { } } }`
	ctx := apc.NewStringContext("<string>", input)
	ctx.AddSkipParser(apc.CastToAny(apc.WhitespaceParser))

	fmt.Printf("Input: %v\n", input)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Node: %v\n", node)
}
