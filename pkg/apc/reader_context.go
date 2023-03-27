package apc

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// ReaderContext[CT] implements Context[CT] by operating with a ReaderWithOrigin[CT].
type ReaderContext[CT any] struct {
	// The reader to use.
	reader ReaderWithOrigin[CT]
	// Buffer used to store read, but unconsumed, elements from reader.
	buffer []CT
	// Buffer used to store read, but unconsumed, element Origins from reader.
	// The elements in this slice must always correspond to the elements in buffer.
	bufferOrigins []Origin
	// The last Origin read from the reader.
	lastOrigin Origin
	// List of parsers to attempt to run, discarding their results if successful.
	skipParsers []Parser[CT, any]
	// Whether or not RunSkipParsers is currently running.
	skipping bool
	// If true, RunSkipParsers will be a no-op. The assumption is that
	// when RunSkipParsers is run, it does not need to be run again until
	// a Consume call.
	skippedSinceLastConsume bool
}

// Returns a *ReaderContext[CT] with the given reader.
func NewReaderContext[CT any](reader ReaderWithOrigin[CT]) *ReaderContext[CT] {
	return &ReaderContext[CT]{
		reader:                  reader,
		buffer:                  make([]CT, 0),
		bufferOrigins:           make([]Origin, 0),
		lastOrigin:              Origin{},
		skipParsers:             make([]Parser[CT, any], 0),
		skipping:                false,
		skippedSinceLastConsume: false,
	}
}

// Returns a *ReaderContext[rune] from an io.RuneReader.
func NewRuneReaderContext(originName string, reader io.RuneReader) *ReaderContext[rune] {
	return NewReaderContext[rune](NewRuneReaderWithOrigin(originName, reader))
}

// Returns a *ReaderContext[byte] from an io.Reader.
func NewByteReaderContext(originName string, reader io.Reader) *ReaderContext[byte] {
	return NewReaderContext[byte](NewByteReaderWithOrigin(originName, reader))
}

// Returns a *ReaderContext[rune] from a string.
func NewStringContext(originName string, data string) *ReaderContext[rune] {
	return NewRuneReaderContext(originName, strings.NewReader(data))
}

// Returns a *ReaderContext[rune] from a file.
func NewFileContext(file *os.File) *ReaderContext[rune] {
	return NewRuneReaderContext(file.Name(), bufio.NewReader(file))
}

// Returns a *ReaderContext[byte] from a file.
func NewBinaryFileContext(file *os.File) *ReaderContext[byte] {
	return NewByteReaderContext(file.Name(), file)
}

// Tries to ensure that num values are in the ctx.buffer. If ErrEOF is reached,
// a nil error is returned here. If another error is reached, that error is returned.
func (ctx *ReaderContext[CT]) maybeEnsureBufferLoaded(num int) error {
	if len(ctx.buffer) >= num {
		return nil
	}
	val, origin, err := ctx.reader.Read()
	for err == nil {
		ctx.buffer = append(ctx.buffer, val)
		ctx.bufferOrigins = append(ctx.bufferOrigins, origin)
		ctx.lastOrigin = origin
		if len(ctx.buffer) >= num {
			break
		}
		val, origin, err = ctx.reader.Read()
	}
	if err != nil && errors.Is(err, ErrEOF) {
		return nil
	}
	return err
}

// Returns a []CT of num elements beginning at offset without consuming
// the elements.
// The offset is a non-negative value relative to the next unconsumed
// element in the input stream.
//
// If the end of input is reached, an EOFError is returned along
// with any peeked elements (which may be less than num elements in length
// if end of input has been reached).
func (ctx *ReaderContext[CT]) Peek(offset int, num int) ([]CT, error) {
	err := ctx.maybeEnsureBufferLoaded(offset + num)
	if err != nil {
		return nil, err
	}
	if len(ctx.buffer) < offset+num {
		return ctx.buffer, ErrEOF
	}
	return ctx.buffer[offset : offset+num], nil
}

// Advances the input stream by num elements, returning the consumed
// elements.
//
// If the end of input is reached, an EOFError is returned along
// with any consumed elements (which may be less than num elements in length
// if end of input has been reached).
func (ctx *ReaderContext[CT]) Consume(num int) ([]CT, error) {
	ctx.skippedSinceLastConsume = false

	err := ctx.maybeEnsureBufferLoaded(num)
	if err != nil {
		return nil, err
	}
	if len(ctx.buffer) < num {
		ret := ctx.buffer[:]
		ctx.buffer = ctx.buffer[:0]
		ctx.bufferOrigins = ctx.bufferOrigins[:0]
		return ret, ErrEOF
	}
	ret := ctx.buffer[:num]
	ctx.buffer = ctx.buffer[num:]
	ctx.bufferOrigins = ctx.bufferOrigins[num:]
	return ret, nil
}

// Returns an Origin representing the next unconsumed element in the
// input stream.
func (ctx *ReaderContext[CT]) GetCurOrigin() Origin {
	ctx.maybeEnsureBufferLoaded(1)
	if len(ctx.bufferOrigins) <= 0 {
		return ctx.lastOrigin
	}
	return ctx.bufferOrigins[0]
}

// Adds the parser to the list of parsers that attempt to run when
// RunSkipParsers is called. If the parser matches, its result will
// be discarded. Duplicate parsers cannot be added.
func (ctx *ReaderContext[CT]) AddSkipParser(parser Parser[CT, any]) {
	for _, p := range ctx.skipParsers {
		if &p == &parser {
			panic("cannot add duplicate skip parser")
		}
	}
	ctx.skipParsers = append(ctx.skipParsers, parser)
}

// Removes the parser from the list of parsers that attempt to run
// when RunSkipParsers is called. If the parser has not been added,
// the function panics.
func (ctx *ReaderContext[CT]) RemoveSkipParser(parser Parser[CT, any]) {
	i := -1
	var p Parser[CT, any]
	for i, p = range ctx.skipParsers {
		if &p == &parser {
			break
		}
	}
	if i == -1 {
		panic("cannot remove non-existent skip parser")
	}
	ctx.skipParsers = append(ctx.skipParsers[:i], ctx.skipParsers[i+1:]...)
}

// Attempts to run any added skip parsers as long as one of the parsers
// successfully matches. The results of any matched parsers is discarded.
// Should only return nil or non-ParseError errors.
func (ctx *ReaderContext[CT]) RunSkipParsers() error {
	if ctx.skipping || ctx.skippedSinceLastConsume {
		return nil
	}

	ctx.skipping = true

	skip := true
	for skip {
		skip = false
		for _, parser := range ctx.skipParsers {
			_, err := parser(ctx)
			if err == nil {
				skip = true
				break
			} else if IsMustReturnParseErr(err) {
				ctx.skipping = false
				return err
			}
		}
	}

	ctx.skippedSinceLastConsume = true
	ctx.skipping = false
	return nil
}
