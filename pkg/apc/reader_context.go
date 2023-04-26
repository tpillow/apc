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
	lastOrigin    Origin
	// List of parsers to attempt to run, discarding their results if successful.
	skipParsers []Parser[CT, any]
	// Whether or not RunSkipParsers is currently running.
	skipping bool
	// If true, RunSkipParsers will be a no-op. The assumption is that
	// when RunSkipParsers is run, it does not need to be run again until
	// a Consume call.
	skippedSinceLastConsume bool
	lookStack               []int
	nameStack               []string
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
		lookStack:               make([]int, 0),
		nameStack:               make([]string, 0),
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
	lookOffset := 0
	if len(ctx.lookStack) > 0 {
		lookOffset = ctx.lookStack[len(ctx.lookStack)-1]
	}

	err := ctx.maybeEnsureBufferLoaded(lookOffset + offset + num)
	if err != nil {
		return nil, err
	}
	buf := ctx.buffer[lookOffset:]
	if len(buf) < offset+num {
		return buf, ErrEOF
	}
	return buf[offset : offset+num], nil
}

// Advances the input stream by num elements, returning the consumed
// elements.
//
// If the end of input is reached, an EOFError is returned along
// with any consumed elements (which may be less than num elements in length
// if end of input has been reached).
func (ctx *ReaderContext[CT]) Consume(num int) ([]CT, error) {
	ctx.skippedSinceLastConsume = false
	lookOffset := 0
	if len(ctx.lookStack) > 0 {
		lookOffset = ctx.lookStack[len(ctx.lookStack)-1]
	}

	err := ctx.maybeEnsureBufferLoaded(lookOffset + num)
	if err != nil {
		return nil, err
	}
	buf := ctx.buffer[lookOffset:]
	if len(buf) < num {
		if len(ctx.lookStack) > 0 {
			ctx.lookStack[len(ctx.lookStack)-1] += len(buf)
		} else {
			ctx.buffer = ctx.buffer[:0]
			ctx.bufferOrigins = ctx.bufferOrigins[:0]
		}
		return buf, ErrEOF
	}
	buf = buf[:num]
	if len(ctx.lookStack) > 0 {
		ctx.lookStack[len(ctx.lookStack)-1] += num
	} else {
		ctx.buffer = ctx.buffer[num:]
		ctx.bufferOrigins = ctx.bufferOrigins[num:]
	}
	return buf, nil
}

// Returns an Origin representing the next unconsumed element in the
// input stream.
func (ctx *ReaderContext[CT]) GetCurOrigin() Origin {
	lookOffset := 0
	if len(ctx.lookStack) > 0 {
		lookOffset = ctx.lookStack[len(ctx.lookStack)-1]
	}

	ctx.maybeEnsureBufferLoaded(lookOffset + 1)
	if len(ctx.bufferOrigins) == 0 {
		return ctx.lastOrigin
	}
	return ctx.bufferOrigins[lookOffset]
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

func (ctx *ReaderContext[CT]) NewLook() {
	if len(ctx.lookStack) == 0 {
		ctx.lookStack = append(ctx.lookStack, 0)
	} else {
		ctx.lookStack = append(ctx.lookStack, ctx.lookStack[len(ctx.lookStack)-1])
	}
}

func (ctx *ReaderContext[CT]) RevertLook() {
	if len(ctx.lookStack) == 0 {
		panic("cannot RevertLook() without a NewLook() on the stack")
	}
	ctx.lookStack = ctx.lookStack[:len(ctx.lookStack)-1]
}

func (ctx *ReaderContext[CT]) CommitLook() error {
	if len(ctx.lookStack) == 0 {
		panic("cannot CommitLook() without a NewLook() on the stack")
	}
	toConsume := ctx.lookStack[len(ctx.lookStack)-1]
	ctx.lookStack = ctx.lookStack[:len(ctx.lookStack)-1]
	if len(ctx.lookStack) == 0 {
		_, err := ctx.Consume(toConsume)
		return err
	} else {
		ctx.lookStack[len(ctx.lookStack)-1] += toConsume - ctx.lookStack[len(ctx.lookStack)-1]
	}
	return nil
}

func (ctx *ReaderContext[CT]) PushName(name string) {
	ctx.nameStack = append(ctx.nameStack, name)
}

func (ctx *ReaderContext[CT]) PopName() {
	if len(ctx.nameStack) == 0 {
		panic("Cannot PopName with nothing on the name stack")
	}
	ctx.nameStack = ctx.nameStack[:len(ctx.nameStack)-1]
}

func (ctx *ReaderContext[CT]) PeekName() string {
	if len(ctx.nameStack) == 0 {
		return "<unknown>"
	}
	return ctx.nameStack[len(ctx.nameStack)-1]
}
