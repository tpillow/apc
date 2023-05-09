package apcgen

import (
	"sort"
	"strings"

	"github.com/tpillow/apc/pkg/apc"
)

func BuildSimpleLexer(
	identifierTokenType apc.TokenType,
	identifierParser apc.Parser[rune, string],
	specialIdentifierTokenTypes []apc.TokenType,
	exactMatchTokenTypes []apc.TokenType,
	providedParsers map[apc.TokenType]apc.Parser[rune, apc.Token],
	skipParsers []apc.Parser[rune, any],
) apc.Parser[rune, apc.Token] {
	sort.Slice(
		exactMatchTokenTypes,
		func(i int, j int) bool {
			return len(exactMatchTokenTypes[i]) > len(exactMatchTokenTypes[j])
		},
	)
	exactMatchTokenLookLen := 0
	if len(exactMatchTokenTypes) > 0 {
		exactMatchTokenLookLen = len(exactMatchTokenTypes[0])
	}

	specialIdentifierTokenTypeMap := map[string]apc.TokenType{}
	for _, tokType := range specialIdentifierTokenTypes {
		specialIdentifierTokenTypeMap[string(tokType)] = tokType
	}

	identLexParser := apc.Map(
		identifierParser,
		func(node string) apc.Token {
			if tokType, ok := specialIdentifierTokenTypeMap[node]; ok {
				return apc.Token{
					Type:  tokType,
					Value: nil,
				}
			}
			return apc.Token{
				Type:  identifierTokenType,
				Value: node,
			}
		},
	)

	parsers := []apc.Parser[rune, apc.Token]{identLexParser}

	if len(exactMatchTokenTypes) > 0 {
		exactMatchLexParser := func(ctx apc.Context[rune]) (apc.Token, error) {
			peekRunes, err := ctx.Peek(0, exactMatchTokenLookLen)
			if err != nil {
				return apc.Token{}, err
			}
			peek := string(peekRunes)
			for _, tokType := range exactMatchTokenTypes {
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

	lexParser := apc.Named(
		"valid token",
		apc.Any(
			parsers...,
		),
	)
	lexParser = wrapWithSkipParsers(lexParser, skipParsers)
	return lexParser
}
