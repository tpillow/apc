package apc

type Context[IT any] interface {
	Peek() IT
	Pop() IT
	IsEof() bool
	CurLocation() Location
	SetLocation(loc Location)
}

func GetContextUnexpected[IT any](ctx Context[IT]) any {
	if ctx.IsEof() {
		return EofString
	}
	return ctx.Peek()
}
