package apc

import "strconv"

// Parses one or more whitespace characters.
var WhitespaceParser = Regex("whitespace", "\\s+")

// Parses a string wrapped in (") characters.
// The result is a string excluding the start and end (").
// The result may contain raw/unescaped (\") and other escape markers.
var DoubleQuotedStringParser = Map(
	Regex("double-quoted string", `"(?:[^"\\]|\\.)*"`),
	func(node string) string {
		return node[1 : len(node)-1]
	})

// Parses a string wrapped in (') characters.
// The result is a string excluding the start and end (').
// The result may contain raw/unescaped (\') and other escape markers.
var SingleQuotedStringParser = Map(
	Regex("single-quoted string", `'(?:[^'\\]|\\.)*'`),
	func(node string) string {
		return node[1 : len(node)-1]
	})

// Parses a C-style identifier.
var IdentifierParser = Regex("identifier", "[a-zA-Z_][a-zA-Z_0-9]*")

// Parses "true" and "false" literals into a boolean.
var BoolParser = OneOf("boolean", Bind(Exact("true"), true), Bind(Exact("false"), false))

// Parses floating point numbers. May be preceded with '+' or '-'.
var FloatParser = Map(
	Regex("float", "[+\\-]?\\d+(\\.\\d+)?"),
	func(node string) float64 {
		val, err := strconv.ParseFloat(node, 64)
		if err != nil {
			panic(err)
		}
		return val
	})

// Parses integer numbers. May be preceded with '+' or '-'.
var IntParser = Map(
	Regex("integer", "[+\\-]?\\d+"),
	func(node string) int64 {
		val, err := strconv.ParseInt(node, 10, 64)
		if err != nil {
			panic(err)
		}
		return val
	})
