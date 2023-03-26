package apc

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

// ReaderContext is a Context that operates off of CT as the
// input stream.
type ReaderContext[CT comparable] struct {
	// Input stream, where index 0 is the next unconsumed CT.
	reader ReaderWithOrigin[CT]
	// Buffer used to store read, but unconsumed, CTs.
	buffer []CT
	// A buffer that matches Origin to each corresponding element in buffer.
	bufferOrigins []Origin
	// The last Origin read from the reader.
	lastOrigin Origin
	// List of parsers to attempt to run, discarding the result.
	skipParsers []Parser[CT, any]
	// Whether or not RunSkipParsers is currently running.
	skipping bool
	// If true, RunSkipParsers will be a no-op. The assumption is that
	// when RunSkipParsers is run, it does not need to be run again until
	// a Consume call.
	skippedSinceLastConsume bool
}

// Returns a *ReaderContext[CT] with the given origin name and CT input stream.
func NewReaderContext[CT comparable](reader ReaderWithOrigin[CT]) *ReaderContext[CT] {
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

// Returns a *ReaderContext[rune] from a string.
func NewStringContext(originName string, data string) *ReaderContext[rune] {
	return NewRuneReaderContext(originName, strings.NewReader(data))
}

// Returns a *ReaderContext[rune] from a file.
func NewFileContext(file *os.File) *ReaderContext[rune] {
	return NewRuneReaderContext(file.Name(), bufio.NewReader(file))
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

func (ctx *ReaderContext[CT]) GetCurOrigin() Origin {
	ctx.maybeEnsureBufferLoaded(1)
	if len(ctx.bufferOrigins) <= 0 {
		return ctx.lastOrigin
	}
	return ctx.bufferOrigins[0]
}

func (ctx *ReaderContext[CT]) AddSkipParser(parser Parser[CT, any]) {
	for _, p := range ctx.skipParsers {
		if &p == &parser {
			panic("cannot add duplicate skip parser")
		}
	}
	ctx.skipParsers = append(ctx.skipParsers, parser)
}

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
