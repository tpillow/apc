package apc

import (
	"strings"
)

// Context holds the current state of the input parsing stream
// and provides methods to peek the input stream, consume it,
// get the current Origin of the input stream, etc.
//
// Also allows for parsers to be added/removed that will skip matched
// input.
type Context[T any] interface {
	// Returns a string of num runes beginning at offset without consuming
	// the runes.
	// The offset is a non-negative value relative to the next unconsumed
	// rune in the input stream.
	//
	// If the end of input is reached, an EOFError is returned along
	// with any peeked runes returned as a string (which may be less
	// than num runes in length if end of input has been reached).
	Peek(offset int, num int) (T, error)
	// Advances the input stream by num runes, returning the consumed
	// runes as a string.
	//
	// If the end of input is reached, an EOFError is returned along
	// with any consumed runes returned as a string (which may be less
	// than num runes in length if end of input has been reached).
	Consume(num int) (T, error)
	// Returns an Origin representing the next unconsumed rune in the
	// input stream.
	GetCurOrigin() Origin

	// Adds the parser to the list of parsers that attempt to run when
	// RunSkipParsers is called. If the parser matches, its result will
	// be discarded. Duplicate parsers cannot be added.
	AddSkipParser(parser Parser[T, any])
	// Removes the parser from the list of parsers that attempt to run
	// when RunSkipParsers is called. If the parser has not been added,
	// the function panics.
	RemoveSkipParser(parser Parser[T, any])
	// Attempts to run any added skip parsers as long as one of the parsers
	// successfully matches. The results of any matched parsers is discarded.
	// Should only return nil or non-ParseError errors.
	RunSkipParsers() error
}

// StringContext is a Context that operates off of []rune as the
// input stream.
type StringContext struct {
	// Input stream, where index 0 is the next unconsumed rune.
	data []rune
	// Current origin of the input stream.
	curOrigin Origin
	// List of parsers to attempt to run, discarding the result.
	skipParsers []Parser[string, any]
	// Whether or not RunSkipParsers is currently running.
	skipping bool
	// If true, RunSkipParsers will be a no-op. The assumption is that
	// when RunSkipParsers is run, it does not need to be run again until
	// a Consume call.
	skippedSinceLastConsume bool
}

// Returns a *StringContext with the given origin name and []rune input stream.
func NewStringContextFromRunes(originName string, data []rune) *StringContext {
	return &StringContext{
		data: data,
		curOrigin: Origin{
			Name:    originName,
			LineNum: 1,
			ColNum:  1,
		},
		skipParsers:             make([]Parser[string, any], 0),
		skipping:                false,
		skippedSinceLastConsume: false,
	}
}

// Returns a *StringContext with the given origin name and string input stream.
func NewStringContext(originName string, data string) *StringContext {
	return NewStringContextFromRunes(originName, []rune(data))
}

func (ctx *StringContext) Peek(offset int, num int) (string, error) {
	var sb strings.Builder

	for i := offset; i < offset+num; i++ {
		if i < len(ctx.data) {
			sb.WriteRune(ctx.data[i])
		} else {
			return sb.String(), &EOFError{}
		}
	}

	return sb.String(), nil
}

func (ctx *StringContext) Consume(num int) (string, error) {
	var sb strings.Builder
	ctx.skippedSinceLastConsume = false

	for i := 0; i < num; i++ {
		if i < len(ctx.data) {
			r := ctx.data[i]
			if r == '\n' {
				ctx.curOrigin.LineNum += 1
				ctx.curOrigin.ColNum = 1
			} else {
				ctx.curOrigin.ColNum += 1
			}
			sb.WriteRune(r)
		} else {
			ctx.data = []rune{}
			return sb.String(), &EOFError{}
		}
	}
	ctx.data = ctx.data[num:]

	return sb.String(), nil
}

func (ctx *StringContext) GetCurOrigin() Origin {
	return ctx.curOrigin
}

func (ctx *StringContext) AddSkipParser(parser Parser[string, any]) {
	for _, p := range ctx.skipParsers {
		if &p == &parser {
			panic("cannot add duplicate skip parser")
		}
	}
	ctx.skipParsers = append(ctx.skipParsers, parser)
}

func (ctx *StringContext) RemoveSkipParser(parser Parser[string, any]) {
	i := -1
	var p Parser[string, any]
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

func (ctx *StringContext) RunSkipParsers() error {
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
