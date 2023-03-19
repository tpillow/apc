package apc

type StringContext struct {
	data      []rune
	curOrigin Origin
}

func NewStringContext(originName string, data string) *StringContext {
	return &StringContext{
		data: []rune(data),
		curOrigin: Origin{
			Name:    originName,
			LineNum: 1,
			ColNum:  1,
		},
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
