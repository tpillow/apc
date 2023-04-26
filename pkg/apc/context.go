package apc

// Value of an invalid Look.
const InvalidLook int = -1

// Context[CT] holds the current state of some input parsing stream
// of type CT, and provides methods to peek the input stream, consume it,
// get the current Origin of the input stream, etc.
//
// Also allows for parsers to be added/removed that will skip matched
// input.
type Context[CT any] interface {
	// Returns a []CT of num elements beginning at offset without consuming
	// the elements.
	// The offset is a non-negative value relative to the next unconsumed
	// element in the input stream.
	//
	// If the end of input is reached, an EOFError is returned along
	// with any peeked elements (which may be less than num elements in length
	// if end of input has been reached).
	Peek(offset int, num int) ([]CT, error)
	// Advances the input stream by num elements, returning the consumed
	// elements.
	//
	// If the end of input is reached, an EOFError is returned along
	// with any consumed elements (which may be less than num elements in length
	// if end of input has been reached).
	Consume(num int) ([]CT, error)
	// Returns an Origin representing the next unconsumed element in the
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
	// Sets the name of all subsequent parsers.
	SetCurName(name string)
	// Gets the current name of parsers.
	GetCurName() string
	// Sets the look value.
	SetLook(val int)
	// Gets the look value.
	GetLook() int
}
