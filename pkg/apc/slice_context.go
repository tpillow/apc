package apc

import "fmt"

type SliceLocation struct {
	Index   int
	LineNum int
	ColNum  int
}

func (loc SliceLocation) String() string {
	return fmt.Sprintf("%d:%d:(#%d)", loc.LineNum, loc.ColNum, loc.Index)
}

func (loc SliceLocation) next(token any) SliceLocation {
	if IsEofToken(token) {
		return loc
	}

	lineNum := loc.LineNum
	colNum := loc.ColNum
	if r, ok := token.(rune); ok {
		if r == '\n' {
			lineNum += 1
			colNum = 0
		} else {
			colNum += 1
		}
	}

	return SliceLocation{
		Index:   loc.Index + 1,
		LineNum: lineNum,
		ColNum:  colNum,
	}
}

type sliceContext struct {
	source   []any
	location SliceLocation
}

type runeSliceContext struct {
	*sliceContext
	sourceAsRuneSlice []rune
}

func NewSliceContext(source []any) Context {
	return &sliceContext{
		source: source,
		location: SliceLocation{
			Index:   0,
			LineNum: 1,
			ColNum:  0,
		},
	}
}

func NewRuneSliceContext(source []rune) Context {
	sourceAny := make([]any, len(source))
	for i := 0; i < len(source); i++ {
		sourceAny[i] = source[i]
	}

	return &runeSliceContext{
		sliceContext:      NewSliceContext(sourceAny).(*sliceContext),
		sourceAsRuneSlice: source,
	}
}

func NewStringContext(source string) Context {
	return NewRuneSliceContext([]rune(source))
}

func (ctx *sliceContext) Peek() any {
	if ctx.location.Index < 0 {
		panic("sliceContext location.Index must be >= 0")
	}
	if ctx.location.Index >= len(ctx.source) {
		return EofValue{}
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

func (ctx *sliceContext) SetLocation(rawLoc Location) {
	loc, ok := rawLoc.(SliceLocation)
	if !ok {
		panic("sliceContext SetLocation must take a SliceLocation type")
	}
	ctx.location = loc
}
