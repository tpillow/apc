package apc

import "fmt"

func basicMatchCallback[IT comparable](what IT, doPop bool) ParserFunc[IT, IT] {
	return func(ctx Context[IT]) (IT, error) {
		if ctx.IsEof() {
			return newT[IT](), ExpErr{
				AtIndex:    ctx.CurIndex(),
				Expected:   fmt.Sprintf("%v", what),
				Unexpected: "EOF",
			}
		}
		val := ctx.Peek()
		if val != what {
			return newT[IT](), ExpErr{
				AtIndex:    ctx.CurIndex(),
				Expected:   fmt.Sprintf("%v", what),
				Unexpected: fmt.Sprintf("%v", val),
			}
		}
		if doPop {
			ctx.Pop()
		}
		return val, nil
	}
}

func Match[IT comparable](what IT) Parser[IT, IT] {
	return NewCallbackParser(basicMatchCallback(what, true))
}

func Test[IT comparable](what IT) Parser[IT, IT] {
	return NewCallbackParser(basicMatchCallback(what, false))
}

func Peek[IT, OT any](parser Parser[IT, OT]) Parser[IT, OT] {
	return NewCallbackParser(func(ctx Context[IT]) (OT, error) {
		indexBefore := ctx.CurIndex()
		result, err := parser.Parse(ctx)
		indexAfter := ctx.CurIndex()
		ctx.Rewind(indexAfter - indexBefore)
		return result, err
	})
}

func Succeed[IT, OT any](val OT) Parser[IT, OT] {
	return NewCallbackParser(func(ctx Context[IT]) (OT, error) {
		return val, nil
	})
}

func Fail[IT any](errMsg string) Parser[IT, IT] {
	return NewCallbackParser(func(ctx Context[IT]) (IT, error) {
		return newT[IT](), CustomErr{
			AtIndex: ctx.CurIndex(),
			Message: errMsg,
		}
	})
}

func Eof[IT any]() Parser[IT, IT] {
	return NewCallbackParser(func(ctx Context[IT]) (IT, error) {
		if !ctx.IsEof() {
			return newT[IT](), ExpErr{
				AtIndex:    ctx.CurIndex(),
				Expected:   "EOF",
				Unexpected: fmt.Sprintf("%v", ctx.Peek()),
			}
		}
		return newT[IT](), nil
	})
}
