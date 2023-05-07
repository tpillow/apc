package main

import (
	"fmt"

	"github.com/tpillow/apc/pkg/apc"
	"github.com/tpillow/apc/pkg/apcgen"
)

type Dir struct {
	Entries []*DirEntry `apc:"'Dir' '{' $.* '}'"`
}

func (d *Dir) String() string {
	return fmt.Sprintf("<Dir Entries=%v>", d.Entries)
}

type DirEntry struct {
	Name   string `apc:"'Entry' '{' $StrParser"`
	Id     int    `apc:"$regex('[0-9]+')"`
	SubDir *Dir   `apc:"$.? '}'"`
}

func (de *DirEntry) String() string {
	return fmt.Sprintf("<DirEntry Name=%v Id=%v SubDir=%v>", de.Name, de.Id, de.SubDir)
}

func main() {
	parser := apcgen.BuildParser[Dir](
		apcgen.DefaultRuneBuildOptions,
		map[string]apc.Parser[rune, any]{
			"StrParser": apc.CastToAny(apc.DoubleQuotedStringParser),
		},
	)

	input := `Dir { Entry { "Name1" 1 } Entry { "Name2" 2 Dir { Entry { "Name3" 3 } } } }`
	ctx := apc.NewStringContext("<string>", input)
	ctx.AddSkipParser(apc.CastToAny(apc.WhitespaceParser))

	fmt.Printf("Input: %v\n", input)
	node, err := apc.Parse[rune](ctx, parser, apc.DefaultParseConfig)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("Node: %v\n", node)
}
