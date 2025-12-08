package apc

type sliceContext[IT any] struct {
	source   []IT
	location Location
}

func NewSliceContext[IT any](sourceName string, source []IT) Context[IT] {
	return &sliceContext[IT]{
		source: source,
		location: Location{
			Name:    sourceName,
			Index:   0,
			LineNum: 1,
			ColNum:  0,
		},
	}
}

func NewStringContext(sourceName string, source string) Context[rune] {
	return NewSliceContext(sourceName, []rune(source))
}

func (ctx *sliceContext[IT]) IsEof() bool {
	return ctx.location.Index >= len(ctx.source)
}

func (ctx *sliceContext[IT]) Peek() IT {
	if ctx.location.Index < 0 {
		panic("sliceContext location.Index must be >= 0")
	}
	if ctx.IsEof() {
		panic("cannot Peek when at EOF")
	}
	return ctx.source[ctx.location.Index]
}

func (ctx *sliceContext[IT]) Pop() IT {
	val := ctx.Peek()
	ctx.location = ctx.location.next(val)
	return val
}

func (ctx *sliceContext[IT]) CurLocation() Location {
	return ctx.location
}

func (ctx *sliceContext[IT]) SetLocation(location Location) {
	ctx.location = location
}
