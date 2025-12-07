package apc

import (
	"fmt"
	"strconv"
)

var (
	Eof                = Describe(Exact(EofValue{}), "<<EOF>>")
	Whitespace         = Describe(Regex(`\s+`), "whitespace")
	OptionalWhitespace = Optional(Whitespace, "")
	AnyChar            = Describe(Regex(`.`), "any character")
	CIdent             = Describe(Regex(`[a-zA-Z_][a-zA-Z_0-9]*`), "C-style identifier")
	Float64            = Describe(Map(Regex(`-?[0-9]+\.[0-9]+`), func(rawValue any) any {
		val, err := strconv.ParseFloat(rawValue.(string), 64)
		if err != nil {
			panic(fmt.Sprintf("unexpected floating-point parse error: %s", err))
		}
		return val
	}), "floating-point number")
	Int64 = Describe(Map(Regex(`-?[0-9]+`), func(rawValue any) any {
		val, err := strconv.ParseInt(rawValue.(string), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("unexpected integer parse error: %s", err))
		}
		return val
	}), "integer number")
	DoubleQuotedString = Describe(RegexGroup(`"([^"\\]*(?:\\.[^"\\]*)*)"`, 1), "double-quoted string")
	SingleQuotedString = Describe(RegexGroup(`'([^'\\]*(?:\\.[^'\\]*)*)'`, 1), "single-quoted string")
)
