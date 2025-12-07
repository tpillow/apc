package apc

type sliceContext struct {
	source           []any
	location         Location
	inBetweenParser  *Parser
	runningPreParser bool
}

type runeSliceContext struct {
	*sliceContext
	sourceAsRuneSlice []rune
}

func NewSliceContext(sourceName string, source []any) Context {
	return &sliceContext{
		source: source,
		location: Location{
			Name:    sourceName,
			Index:   0,
			LineNum: 1,
			ColNum:  0,
		},
		inBetweenParser:  nil,
		runningPreParser: false,
	}
}

func NewRuneSliceContext(sourceName string, source []rune) Context {
	sourceAny := make([]any, len(source))
	for i := 0; i < len(source); i++ {
		sourceAny[i] = source[i]
	}

	return &runeSliceContext{
		sliceContext:      NewSliceContext(sourceName, sourceAny).(*sliceContext),
		sourceAsRuneSlice: source,
	}
}

func NewStringContext(sourceName string, source string) Context {
	return NewRuneSliceContext(sourceName, []rune(source))
}

func (ctx *sliceContext) IsEof() bool {
	return ctx.location.Index >= len(ctx.source)
}

func (ctx *sliceContext) Peek() any {
	if ctx.location.Index < 0 {
		panic("sliceContext location.Index must be >= 0")
	}
	if ctx.IsEof() {
		panic("cannot Peek when at EOF")
	}
	return ctx.source[ctx.location.Index]
}

func (ctx *sliceContext) Pop() any {
	val := ctx.Peek()
	ctx.location = ctx.location.next(val)
	return val
}

func (ctx *sliceContext) CurLocation() Location {
	return ctx.location
}

func (ctx *sliceContext) SetLocation(location Location) {
	ctx.location = location
}

func (ctx *sliceContext) SetPreParser(parser *Parser) {
	ctx.inBetweenParser = parser
}

func (ctx *sliceContext) GetPreParser() *Parser {
	return ctx.inBetweenParser
}

func (ctx *sliceContext) SetRunningPreParser(running bool) {
	ctx.runningPreParser = running
}

func (ctx *sliceContext) IsRunningPreParser() bool {
	return ctx.runningPreParser
}
