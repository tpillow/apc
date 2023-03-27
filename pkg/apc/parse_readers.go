package apc

import (
	"io"
)

// A reader that provides an Origin per read element of type T.
type ReaderWithOrigin[T any] interface {
	// Returns the next element in a stream of elements of type T,
	// along with the Origin associated with the element.
	// If an error occurs or if no element is available, an error is returned.
	Read() (T, Origin, error)
}

// Implements ReaderWithOrigin[T] by calling the provided parser with the provided
// ctx each time Read is called.
type ParseReader[CT, T any] struct {
	ctx    Context[CT]
	parser Parser[CT, T]
}

// Returns a *ParseReader[CT, T] with the provided ctx and parser.
func NewParseReader[CT, T any](ctx Context[CT], parser Parser[CT, T]) *ParseReader[CT, T] {
	return &ParseReader[CT, T]{
		ctx:    ctx,
		parser: parser,
	}
}

// Calls the parser with the corresponding ctx, returning the result and Origin of the result.
// If an error occurs or if no element is available, an error is returned.
func (r *ParseReader[CT, T]) Read() (T, Origin, error) {
	origin := r.ctx.GetCurOrigin()
	val, err := r.parser(r.ctx)
	return val, origin, err
}

// Implements ReaderWithOrigin[rune] by calling reader.ReadRune.
type RuneReaderWithOrigin struct {
	reader    io.RuneReader
	curOrigin Origin
}

// Returns a *RuneReaderWithOrigin with the provided origin name and reader.
func NewRuneReaderWithOrigin(originName string, reader io.RuneReader) *RuneReaderWithOrigin {
	return &RuneReaderWithOrigin{
		reader: reader,
		curOrigin: Origin{
			Name:    originName,
			LineNum: 1,
			ColNum:  1,
		},
	}
}

// Calls reader.ReadRune, returning the resulting rune and Origin of the rune.
// If an error occurs or if no rune is available, an error is returned.
func (r *RuneReaderWithOrigin) Read() (rune, Origin, error) {
	rn, _, err := r.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return rune(-1), r.curOrigin, ErrEOF
		}
		return rune(-1), r.curOrigin, err
	}

	origin := r.curOrigin
	if rn == '\n' {
		r.curOrigin.LineNum += 1
		r.curOrigin.ColNum = 1
	} else {
		r.curOrigin.ColNum += 1
	}

	return rn, origin, nil
}
