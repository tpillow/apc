package apc

import "strconv"

var WhitespaceParser = Regex("whitespace", "\\s+")

var DoubleQuotedStringParser = Regex("double-quoted string", `"(?:[^"\\]|\\.)*"`)

var SingleQuotedStringParser = Regex("single-quoted string", `'(?:[^'\\]|\\.)*'`)

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
