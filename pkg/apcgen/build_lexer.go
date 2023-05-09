package apcgen

import (
	"sort"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

type SimpleLexerBuildOptions struct {
	identifierTokenType         apc.TokenType
	identifierParser            apc.Parser[rune, string]
	specialIdentifierTokenTypes []apc.TokenType
	exactMatchTokenTypes        []apc.TokenType
	providedParsers             []apc.Parser[rune, apc.Token]
	skipParsers                 []apc.Parser[rune, any]
}

func BuildSimpleLexer(
	opts SimpleLexerBuildOptions,
) apc.Parser[rune, apc.Token] {
	sort.Slice(
		opts.exactMatchTokenTypes,
		func(i int, j int) bool {
			return len(opts.exactMatchTokenTypes[i]) > len(opts.exactMatchTokenTypes[j])
		},
	)
	exactMatchTokenLookLen := 0
	if len(opts.exactMatchTokenTypes) > 0 {
		exactMatchTokenLookLen = len(opts.exactMatchTokenTypes[0])
	}

	specialIdentifierTokenTypeMap := map[string]apc.TokenType{}
	for _, tokType := range opts.specialIdentifierTokenTypes {
		specialIdentifierTokenTypeMap[string(tokType)] = tokType
	}

	identLexParser := apc.Map(
		opts.identifierParser,
		func(node string) apc.Token {
			if tokType, ok := specialIdentifierTokenTypeMap[node]; ok {
				return apc.Token{
					Type:  tokType,
					Value: nil,
				}
			}
			return apc.Token{
				Type:  opts.identifierTokenType,
				Value: node,
			}
		},
	)

	parsers := []apc.Parser[rune, apc.Token]{identLexParser}

	if len(opts.exactMatchTokenTypes) > 0 {
		exactMatchLexParser := func(ctx apc.Context[rune]) (apc.Token, error) {
			peekRunes, err := ctx.Peek(0, exactMatchTokenLookLen)
			if err != nil {
				return apc.Token{}, err
			}
			peek := string(peekRunes)
			for _, tokType := range opts.exactMatchTokenTypes {
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

	parsers = append(parsers, opts.providedParsers...)

	lexParser := apc.Named(
		"valid token",
		apc.Any(
			parsers...,
		),
	)
	lexParser = wrapWithSkipParsers(lexParser, opts.skipParsers)
	return lexParser
}
