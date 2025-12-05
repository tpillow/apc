package apc

import (
	"fmt"
	"strconv"
)

var (
	Whitespace         = Regex(`\s+`).Describe("whitespace")
	OptionalWhitespace = Whitespace.Optional("")
	Float              = Regex(`-?[0-9]+\.[0-9]+`).Map(func(rawValue any) any {
		val, err := strconv.ParseFloat(rawValue.(string), 64)
		if err != nil {
			panic(fmt.Sprintf("unexpected floating-point parse error: %s", err))
		}
		return val
	}).Describe("floating-point number")
	Int = Regex(`-?[0-9]+`).Map(func(rawValue any) any {
		val, err := strconv.ParseInt(rawValue.(string), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("unexpected integer parse error: %s", err))
		}
		return val
	}).Describe("integer number")
	DoubleQuotedString = RegexGroup(`"([^"\\]*(?:\\.[^"\\]*)*)"`, 1).Describe("double-quoted string")
	SingleQuotedString = RegexGroup(`'([^'\\]*(?:\\.[^'\\]*)*)'`, 1).Describe("single-quoted string")
)
