package apc

import "fmt"

func SeqAt[IT, OT any](index int, parsers ...Parser[IT, OT]) Parser[IT, OT] {
	return Named(fmt.Sprintf("sequence at index %d", index), Index(index, Seq(parsers...)))
}

func Head[IT, OT any](parsers ...Parser[IT, OT]) Parser[IT, OT] {
	return Named("head of", SeqAt(0, parsers...))
}

func Tail[IT, OT any](parsers ...Parser[IT, OT]) Parser[IT, OT] {
	return Named("tail of", SeqAt(len(parsers)-1, parsers...))
}

func Many[IT, OT any](parser Parser[IT, OT]) Parser[IT, []OT] {
	return Named("many", Range(parser, 0, -1))
}

func Maybe[IT, OT any](parser Parser[IT, OT], defaultValue OT) Parser[IT, OT] {
	return Named("maybe", NewCallbackParser(func(ctx Context[IT]) (OT, error) {
		results, err := Range(parser, 0, 1).Parse(ctx)
		if err != nil {
			return newT[OT](), err
		}
		if len(results) > 1 {
			panic("unreachable")
		}
		if len(results) == 0 {
			return defaultValue, nil
		}
		return results[0], nil
	}))
}
