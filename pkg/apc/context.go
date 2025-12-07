package apc

type Context interface {
	Peek() any
	Pop() any
	IsEof() bool
	CurLocation() Location
	SetLocation(loc Location)
	// TODO: refactor out below to not be required in interface...
	SetPreParser(parser *Parser)
	GetPreParser() *Parser
	SetRunningPreParser(running bool)
	IsRunningPreParser() bool
}

func GetContextUnexpected(ctx Context) any {
	if ctx.IsEof() {
		return EofString
	}
	return ctx.Peek()
}
