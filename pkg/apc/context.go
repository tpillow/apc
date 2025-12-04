package apc

type Location interface{}

type EofToken struct{}

func IsEofToken(thing any) bool {
	_, ok := thing.(EofToken)
	return ok
}

type Context interface {
	Peek() any
	Pop() any
	CurLocation() Location
	SetLocation(loc Location)
}

func ContextIsEof(ctx Context) bool {
	return IsEofToken(ctx.Peek())
}
