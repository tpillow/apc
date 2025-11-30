package apc

func Whitespace() Parser[rune, string] {
	return Named("whitespace", Regex(`\s+`))
}

func AlphaWord() Parser[rune, string] {
	return Named("alphabetic word", Regex(`[a-zA-Z]+`))
}

func AlphaNumericWord() Parser[rune, string] {
	return Named("alpha-numeric word", Regex(`[a-zA-Z0-9]+`))
}

func IdentWord() Parser[rune, string] {
	return Named("identifier (starting with a letter or '_', then alphanumeric or '_')", Regex(`[a-zA-Z_][a-zA-Z_0-9]*`))
}
