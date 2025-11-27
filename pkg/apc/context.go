package apc

type Context[IT any] interface {
	CurIndex() int
	TotalLength() int
	Peek() IT
	PeekN(count int) []IT
	Pop() IT
	PopN(count int) []IT
	Rewind(count int)
	IsEof() bool
}

type sliceContext[IT any] struct {
	source   []IT
	curIndex int
}

func NewSliceContext[IT any](source []IT) Context[IT] {
	return &sliceContext[IT]{
		source:   source,
		curIndex: 0,
	}
}

func (ctx *sliceContext[IT]) CurIndex() int {
	return ctx.curIndex
}

func (ctx *sliceContext[IT]) TotalLength() int {
	return len(ctx.source)
}

func (ctx *sliceContext[IT]) Peek() IT {
	if ctx.IsEof() {
		panic("cannot Peek or Pop when at EOF")
	}
	return ctx.source[ctx.curIndex]
}

func (ctx *sliceContext[IT]) PeekN(count int) []IT {
	if count == 0 {
		panic("cannot PeekN or PopN with count == 0")
	}
	if count < 0 {
		return ctx.source[ctx.CurIndex():]
	}
	if ctx.CurIndex()+count >= ctx.TotalLength() {
		panic("cannot PeekN or PopN with a count that goes past EOF")
	}
	return ctx.source[ctx.CurIndex() : ctx.CurIndex()+count]
}

func (ctx *sliceContext[IT]) PopN(count int) []IT {
	result := ctx.PeekN(count)
	if count < 0 {
		ctx.curIndex = ctx.TotalLength()
	} else {
		ctx.curIndex += count
	}
	return result
}

func (ctx *sliceContext[IT]) Pop() IT {
	val := ctx.Peek()
	ctx.curIndex += 1
	return val
}

func (ctx *sliceContext[IT]) Rewind(count int) {
	if count <= 0 {
		panic("cannot Rewind with count <= 0")
	}
	newIndex := ctx.curIndex - count
	if newIndex < 0 {
		panic("cannot Rewind past index 0")
	}
	ctx.curIndex = newIndex
}

func (ctx *sliceContext[IT]) IsEof() bool {
	return ctx.CurIndex() < ctx.TotalLength()
}
