package apcgen

import (
	"sort"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

type SimpleLexerBuildOptions struct {
	IdentifierTokenType         apc.TokenType
	IdentifierParser            apc.Parser[rune, string]
	SpecialIdentifierTokenTypes []apc.TokenType
	ExactMatchTokenTypes        []apc.TokenType
	ProvidedParsers             []apc.Parser[rune, apc.Token]
	SkipParsers                 []apc.Parser[rune, any]
}

func BuildSimpleLexer(
	opts SimpleLexerBuildOptions,
) apc.Parser[rune, apc.Token] {
	sort.Slice(
		opts.ExactMatchTokenTypes,
		func(i int, j int) bool {
			return len(opts.ExactMatchTokenTypes[i]) > len(opts.ExactMatchTokenTypes[j])
		},
	)
	exactMatchTokenLookLen := 0
	if len(opts.ExactMatchTokenTypes) > 0 {
		exactMatchTokenLookLen = len(opts.ExactMatchTokenTypes[0])
	}

	specialIdentifierTokenTypeMap := map[string]apc.TokenType{}
	for _, tokType := range opts.SpecialIdentifierTokenTypes {
		specialIdentifierTokenTypeMap[string(tokType)] = tokType
	}

	identLexParser := apc.Map(
		opts.IdentifierParser,
		func(node string) apc.Token {
			if tokType, ok := specialIdentifierTokenTypeMap[node]; ok {
				return apc.Token{
					Type:  tokType,
					Value: nil,
				}
			}
			return apc.Token{
				Type:  opts.IdentifierTokenType,
				Value: node,
			}
		},
	)

	parsers := []apc.Parser[rune, apc.Token]{identLexParser}

	if len(opts.ExactMatchTokenTypes) > 0 {
		exactMatchLexParser := func(ctx apc.Context[rune]) (apc.Token, error) {
			peekRunes, err := ctx.Peek(0, exactMatchTokenLookLen)
			if err != nil {
				return apc.Token{}, err
			}
			peek := string(peekRunes)
			for _, tokType := range opts.ExactMatchTokenTypes {
				tokTypeStr := string(tokType)
				if strings.HasPrefix(peek, tokTypeStr) {
					_, err := ctx.Consume(len(tokTypeStr))
					if err != nil {
						return apc.Token{}, err
					}
					return apc.Token{
						Type:  tokType,
						Value: nil,
					}, nil
				}
			}
			return apc.Token{}, apc.ParseErrExpectedButGotNext(ctx, "valid token", nil)
		}

		parsers = append(parsers, exactMatchLexParser)
	}

	parsers = append(parsers, opts.ProvidedParsers...)

	lexParser := apc.Named(
		"valid token",
		apc.Any(
			parsers...,
		),
	)
	lexParser = wrapWithSkipParsers(lexParser, opts.SkipParsers)
	return lexParser
}
