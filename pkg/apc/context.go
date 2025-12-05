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
	// TODO: refactor out below to not be required in interface...
	SetPreParser(parser *Parser)
	GetPreParser() *Parser
	SetRunningPreParser(running bool)
	IsRunningPreParser() bool
}

func ContextIsEof(ctx Context) bool {
	return IsEofValue(ctx.Peek())
}
