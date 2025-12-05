package apc

import (
	"fmt"
	"strconv"
)

var (
	Eof                = Exact(EofToken{}).Describe("<<EOF>>")
	Whitespace         = Regex(`\s+`).Describe("whitespace")
	OptionalWhitespace = Whitespace.Optional("")
	AnyChar            = Regex(`.`).Describe("any character")
	CIdent             = Regex(`[a-zA-Z_][a-zA-Z_0-9]*`).Describe("C-style identifier")
	Float64            = Regex(`-?[0-9]+\.[0-9]+`).Map(func(rawValue any) any {
		val, err := strconv.ParseFloat(rawValue.(string), 64)
		if err != nil {
			panic(fmt.Sprintf("unexpected floating-point parse error: %s", err))
		}
		return val
	}).Describe("floating-point number")
	Int64 = Regex(`-?[0-9]+`).Map(func(rawValue any) any {
		val, err := strconv.ParseInt(rawValue.(string), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("unexpected integer parse error: %s", err))
		}
		return val
	}).Describe("integer number")
	DoubleQuotedString = RegexGroup(`"([^"\\]*(?:\\.[^"\\]*)*)"`, 1).Describe("double-quoted string")
	SingleQuotedString = RegexGroup(`'([^'\\]*(?:\\.[^'\\]*)*)'`, 1).Describe("single-quoted string")
)
