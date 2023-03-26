package apc

// Context holds the current state of the input parsing stream
// and provides methods to peek the input stream, consume it,
// get the current Origin of the input stream, etc.
//
// Also allows for parsers to be added/removed that will skip matched
// input.
type Context[CT comparable] interface {
	// Returns a string of num runes beginning at offset without consuming
	// the runes.
	// The offset is a non-negative value relative to the next unconsumed
	// rune in the input stream.
	//
	// If the end of input is reached, an EOFError is returned along
	// with any peeked runes returned as a string (which may be less
	// than num runes in length if end of input has been reached).
	Peek(offset int, num int) ([]CT, error)
	// Advances the input stream by num runes, returning the consumed
	// runes as a string.
	//
	// If the end of input is reached, an EOFError is returned along
	// with any consumed runes returned as a string (which may be less
	// than num runes in length if end of input has been reached).
	Consume(num int) ([]CT, error)
	// Returns an Origin representing the next unconsumed rune in the
	// input stream.
	GetCurOrigin() Origin

	// Adds the parser to the list of parsers that attempt to run when
	// RunSkipParsers is called. If the parser matches, its result will
	// be discarded. Duplicate parsers cannot be added.
	AddSkipParser(parser Parser[CT, any])
	// Removes the parser from the list of parsers that attempt to run
	// when RunSkipParsers is called. If the parser has not been added,
	// the function panics.
	RemoveSkipParser(parser Parser[CT, any])
	// Attempts to run any added skip parsers as long as one of the parsers
	// successfully matches. The results of any matched parsers is discarded.
	// Should only return nil or non-ParseError errors.
	RunSkipParsers() error
}
