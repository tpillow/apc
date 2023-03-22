package apc

type Parser[T any] func(ctx Context) (T, error)

type ParseConfig struct {
	MustParseToEOF bool
}

func Parse[T any](ctx Context, parser Parser[T], parseConfig ParseConfig) (T, error) {
	ctx.RunSkipParsers()
	node, err := parser(ctx)
	if err != nil {
		return zeroVal[T](), err
	}

	if parseConfig.MustParseToEOF {
		ctx.RunSkipParsers()
		if _, err := ctx.Peek(0, 1); err == nil {
			return zeroVal[T](), ParseErrExpectedButGotNext(ctx, "EOF", nil)
		}
	}

	return node, nil
}
