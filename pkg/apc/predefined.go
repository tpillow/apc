package apc

import "strconv"

// Parses one or more whitespace characters and returns the string result.
var WhitespaceParser = Named("whitespace", Regex("\\s+"))

// Parses a string wrapped in (") characters.
// The result is a string excluding the start and end (").
// The result may contain raw/unescaped (\") and other escape markers.
var DoubleQuotedStringParser = Named("double-quoted string", Map(
	Regex(`"(?:[^"\\]|\\.)*"`),
	func(node string) string {
		return node[1 : len(node)-1]
	}))

// Parses a string wrapped in (') characters.
// The result is a string excluding the start and end (').
// The result may contain raw/unescaped (\') and other escape markers.
var SingleQuotedStringParser = Named("single-quoted string", Map(
	Regex(`'(?:[^'\\]|\\.)*'`),
	func(node string) string {
		return node[1 : len(node)-1]
	}))

// Parses a C-style identifier and returns the string result.
var IdentifierParser = Named("identifier", Regex("[a-zA-Z_][a-zA-Z_0-9]*"))

// Parses "true" and "false" literals into a boolean and returns
// the boolean result.
var BoolParser = Named("boolean", Any(Bind(ExactStr("true"), true), Bind(ExactStr("false"), false)))

// Parses floating point numbers and returns a float64 result.
// May be preceded with '+' or '-'.
var FloatParser = Named("float",
	Map(
		Regex("[+\\-]?\\d+(\\.\\d+)?"),
		func(node string) float64 {
			val, err := strconv.ParseFloat(node, 64)
			if err != nil {
				panic(err)
			}
			return val
		}))

// Parses integer numbers and returns an int64 result.
// May be preceded with '+' or '-'.
var IntParser = Named("integer",
	Map(
		Regex("[+\\-]?\\d+"),
		func(node string) int64 {
			val, err := strconv.ParseInt(node, 10, 64)
			if err != nil {
				panic(err)
			}
			return val
		}))
