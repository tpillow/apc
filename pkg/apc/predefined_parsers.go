package apc

import "strconv"

var WhitespaceParser = Regex("whitespace", "\\s+")

var DoubleQuotedStringParser = Map(
	Regex("double-quoted string", `"(?:[^"\\]|\\.)*"`),
	func(node string) string {
		return node[1 : len(node)-1]
	})

var SingleQuotedStringParser = Map(
	Regex("single-quoted string", `'(?:[^'\\]|\\.)*'`),
	func(node string) string {
		return node[1 : len(node)-1]
	})

var IdentifierParser = Regex("identifier", "[a-zA-Z_][a-zA-Z_0-9]*")

var BoolParser = Any("boolean", Bind(Exact("true"), true), Bind(Exact("false"), false))

var FloatParser = Map(
	Regex("float", "[+\\-]?\\d+(\\.\\d+)?"),
	func(node string) float64 {
		val, err := strconv.ParseFloat(node, 64)
		if err != nil {
			panic(err)
		}
		return val
	})

var IntParser = Map(
	Regex("integer", "[+\\-]?\\d+"),
	func(node string) int64 {
		val, err := strconv.ParseInt(node, 10, 64)
		if err != nil {
			panic(err)
		}
		return val
	})
