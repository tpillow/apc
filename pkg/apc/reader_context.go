package apc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
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
	// Current parser name.
	curParserName string
	// Current look offset value.
	lookOffset int
	// Current debug indentation.
	debugIndentation string
	// Whether or not to enable debugging.
	DebugParsers bool
}

// Returns a *ReaderContext[CT] with the given reader.
func NewReaderContext[CT any](reader ReaderWithOrigin[CT]) *ReaderContext[CT] {
	return &ReaderContext[CT]{
		reader:           reader,
		buffer:           make([]CT, 0),
		bufferOrigins:    make([]Origin, 0),
		lastOrigin:       Origin{},
		skipParsers:      make([]Parser[CT, any], 0),
		skipping:         false,
		curParserName:    "<unknown>",
		lookOffset:       InvalidLookOffset,
		debugIndentation: "",
		DebugParsers:     false,
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
		var tmpRune rune
		if reflect.TypeOf(val) == reflect.TypeOf(tmpRune) {
			maybeLog(DebugPrintReaderContext, "ReaderContext %p appended to buffer: %c", ctx, val)
		} else {
			maybeLog(DebugPrintReaderContext, "ReaderContext %p appended to buffer: %v", ctx, val)
		}
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
	if ctx.lookOffset != InvalidLookOffset {
		lookOffset = ctx.lookOffset
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
	lookOffset := 0
	if ctx.lookOffset != InvalidLookOffset {
		lookOffset = ctx.lookOffset
	}

	err := ctx.maybeEnsureBufferLoaded(lookOffset + num)
	if err != nil {
		return nil, err
	}
	buf := ctx.buffer[lookOffset:]
	if len(buf) < num {
		if ctx.lookOffset != InvalidLookOffset {
			ctx.lookOffset += len(buf)
		} else {
			ctx.buffer = ctx.buffer[:0]
			ctx.bufferOrigins = ctx.bufferOrigins[:0]
		}
		return buf, ErrEOF
	}
	buf = buf[:num]
	if ctx.lookOffset != InvalidLookOffset {
		ctx.lookOffset += num
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
	if ctx.lookOffset != InvalidLookOffset {
		lookOffset = ctx.lookOffset
	}

	ctx.maybeEnsureBufferLoaded(lookOffset + 1)
	if len(ctx.bufferOrigins) <= lookOffset {
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
	if ctx.skipping {
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

	ctx.skipping = false
	return nil
}

// Sets the look value.
func (ctx *ReaderContext[CT]) SetLookOffset(val int) {
	ctx.lookOffset = val
}

// Gets the look value.
func (ctx *ReaderContext[CT]) GetLookOffset() int {
	return ctx.lookOffset
}

// Sets the name of all subsequent parsers.
func (ctx *ReaderContext[CT]) SetCurParserName(name string) {
	ctx.curParserName = name
}

// Gets the current name of parsers.
func (ctx *ReaderContext[CT]) GetCurParserName() string {
	return ctx.curParserName
}

// TODO: document
func (ctx *ReaderContext[CT]) DebugStart(format string, formatArgs ...interface{}) {
	if !ctx.DebugParsers {
		return
	}
	fmt.Printf("%vSTART: %v (in %v) @ %v\n", ctx.debugIndentation, fmt.Sprintf(format, formatArgs...), ctx.GetCurParserName(), ctx.GetCurOrigin())
	ctx.debugIndentation += "  "
}

// TODO: document
func (ctx *ReaderContext[CT]) DebugPrint(format string, formatArgs ...interface{}) {
	if !ctx.DebugParsers {
		return
	}
	fmt.Printf("%vDEBUG: %v (in %v) @ %v\n", ctx.debugIndentation, fmt.Sprintf(format, formatArgs...), ctx.GetCurParserName(), ctx.GetCurOrigin())
}

// TODO: document
func (ctx *ReaderContext[CT]) DebugEnd(format string, formatArgs ...interface{}) {
	if !ctx.DebugParsers {
		return
	}
	ctx.debugIndentation = ctx.debugIndentation[:len(ctx.debugIndentation)-2]
	fmt.Printf("%vEND: %v @ %v\n", ctx.debugIndentation, fmt.Sprintf(format, formatArgs...), ctx.GetCurOrigin())
}
