package apc

type Location interface{}

type EofValue struct{}

func (EofValue) String() string {
	return "<<EOF>>"
}

func IsEofValue(thing any) bool {
	_, ok := thing.(EofValue)
	return ok
}

type Context interface {
	Peek() any
	Pop() any
	CurLocation() Location
	SetLocation(loc Location)
}

func ContextIsEof(ctx Context) bool {
	return IsEofValue(ctx.Peek())
}
