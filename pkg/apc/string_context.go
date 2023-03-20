package apc

import "fmt"

type StringContext struct {
	data        []rune
	curOrigin   Origin
	skipParsers map[Parser]bool
}

func NewStringContext(originName string, data string) *StringContext {
	return &StringContext{
		data: []rune(data),
		curOrigin: Origin{
			Name:    originName,
			LineNum: 1,
			ColNum:  1,
		},
		skipParsers: make(map[Parser]bool),
	}
}

func (ctx *StringContext) PeekRune(offset int) (rune, error) {
	if len(ctx.data) <= offset {
		return 0, &EOFError{}
	}
	return ctx.data[offset], nil
}

func (ctx *StringContext) ConsumeRune() (rune, error) {
	if len(ctx.data) <= 0 {
		return 0, &EOFError{}
	}
	r := ctx.data[0]
	ctx.data = ctx.data[1:]
	switch r {
	case '\n':
		ctx.curOrigin.LineNum += 1
		ctx.curOrigin.ColNum = 1
	default:
		ctx.curOrigin.ColNum += 1
	}
	return r, nil
}

func (ctx *StringContext) GetOrigin() Origin {
	return ctx.curOrigin
}

func (ctx *StringContext) AddSkipParser(parser Parser) {
	if _, has := ctx.skipParsers[parser]; has {
		panic(fmt.Errorf("cannot add duplicate parser %v to skip", parser))
	}
	ctx.skipParsers[parser] = true
}

func (ctx *StringContext) RemoveSkipParser(parser Parser) {
	if _, has := ctx.skipParsers[parser]; !has {
		panic(fmt.Errorf("cannot remove un-added skip parser %v", parser))
	}
	delete(ctx.skipParsers, parser)
}

func (ctx *StringContext) ProcessSkips() {
	skipped := true
	for skipped {
		skipped = false
		for p := range ctx.skipParsers {
			if _, err := p.Parse(ctx); err == nil {
				skipped = true
				break
			}
		}
	}
}
