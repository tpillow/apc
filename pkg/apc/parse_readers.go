package apc

import (
	"io"
)

type ReaderWithOrigin[T any] interface {
	Read() (T, Origin, error)
}

type ParseReader[CT comparable, T any] struct {
	ctx    Context[CT]
	parser Parser[CT, T]
}

func NewParseReader[CT comparable, T any](ctx Context[CT], parser Parser[CT, T]) *ParseReader[CT, T] {
	return &ParseReader[CT, T]{
		ctx:    ctx,
		parser: parser,
	}
}

func (r *ParseReader[CT, T]) Read() (T, Origin, error) {
	origin := r.ctx.GetCurOrigin()
	val, err := r.parser(r.ctx)
	return val, origin, err
}

type RuneReaderWithOrigin struct {
	reader    io.RuneReader
	curOrigin Origin
}

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

type Lexer[CT comparable, T any] struct {
	ctx    Context[CT]
	parser Parser[CT, T]
}

func NewLexer[CT comparable, T any](ctx Context[CT], parser Parser[CT, T]) *Lexer[CT, T] {
	return &Lexer[CT, T]{
		ctx:    ctx,
		parser: parser,
	}
}

func (r *Lexer[CT, T]) Read() (T, Origin, error) {
	val, err := r.parser(r.ctx)
	if err != nil {
		return zeroVal[T](), r.ctx.GetCurOrigin(), err
	}
	return val, r.ctx.GetCurOrigin(), nil
}
