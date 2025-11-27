package apc

type ParserFunc[IT, OT any] func(ctx Context[IT]) (OT, error)

type Parser[IT, OT any] interface {
	Parse(ctx Context[IT]) (OT, error)
}

type callbackParser[IT, OT any] struct {
	callback ParserFunc[IT, OT]
}

func (p callbackParser[IT, OT]) Parse(ctx Context[IT]) (OT, error) {
	return p.callback(ctx)
}

func NewCallbackParser[IT, OT any](callback ParserFunc[IT, OT]) Parser[IT, OT] {
	return callbackParser[IT, OT]{
		callback: callback,
	}
}

type ParserRef[IT, OT any] struct {
	Parser Parser[IT, OT]
}

func (p ParserRef[IT, OT]) Parse(ctx Context[IT]) (OT, error) {
	if p.Parser == nil {
		panic("cannot call Parse on a ParserRef that does not have its Parser attribute set")
	}
	return p.Parser.Parse(ctx)
}
