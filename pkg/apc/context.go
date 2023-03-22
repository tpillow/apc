package apc

import (
	"strings"
)

type Context interface {
	Peek(offset int, num int) (string, error)
	Consume(num int) (string, error)
	GetCurOrigin() Origin

	AddSkipParser(parser Parser[any])
	RemoveSkipParser(parser Parser[any])
	RunSkipParsers()
}

type StringContext struct {
	data                    []rune
	curOrigin               Origin
	skipParsers             []Parser[any]
	skipping                bool
	skippedSinceLastConsume bool
}

func NewStringContextFromRunes(originName string, data []rune) *StringContext {
	return &StringContext{
		data: data,
		curOrigin: Origin{
			Name:    originName,
			LineNum: 1,
			ColNum:  1,
		},
		skipParsers:             make([]Parser[any], 0),
		skipping:                false,
		skippedSinceLastConsume: false,
	}
}

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

func (ctx *StringContext) AddSkipParser(parser Parser[any]) {
	for _, p := range ctx.skipParsers {
		if &p == &parser {
			panic("cannot add duplicate skip parser")
		}
	}
	ctx.skipParsers = append(ctx.skipParsers, parser)
}

func (ctx *StringContext) RemoveSkipParser(parser Parser[any]) {
	i := -1
	var p Parser[any]
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

func (ctx *StringContext) RunSkipParsers() {
	if ctx.skipping || ctx.skippedSinceLastConsume {
		return
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
			}
		}
	}
	ctx.skippedSinceLastConsume = true
	ctx.skipping = false
}
