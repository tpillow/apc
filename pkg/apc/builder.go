package apc

type Builder[IT, OT any] struct {
	parser *Parser[IT, OT]
}

func NewBuilder[IT, OT any](initialParser *Parser[IT, OT]) *Builder[IT, OT] {
	// TODO: no pointer type
	return &Builder[IT, OT]{
		parser: initialParser,
	}
}

func (builder *Builder[IT, OT]) Build() *Parser[IT, OT] {
	return builder.parser
}

func (builder *Builder[IT, OT]) Describe(description string) *Builder[IT, OT] {
	builder.parser = Describe(builder.parser, description)
	return builder
}

func (builder *Builder[IT, OT]) Peek() *Builder[IT, OT] {
	builder.parser = Peek(builder.parser)
	return builder
}

func (builder *Builder[IT, OT]) Index(index int) *Builder[IT, OT] {
	builder.parser = Index[IT, OT](builder.parser, index)
	return builder
}

func (builder *Builder[IT, OT]) ConcatSlices(parser *Parser) *Builder {
	builder.parser = ConcatSlices(builder.parser, parser)
	return builder
}

func (builder *Builder[IT, OT]) Times(min int, max int) *Builder {
	builder.parser = Times(builder.parser, min, max)
	return builder
}

func (builder *Builder[IT, OT]) Optional(defaultValue any) *Builder {
	builder.parser = Optional(builder.parser, defaultValue)
	return builder
}

func (builder *Builder[IT, OT]) Head() *Builder {
	builder.parser = Head(builder.parser)
	return builder
}

func (builder *Builder[IT, OT]) Tail() *Builder {
	builder.parser = Tail(builder.parser)
	return builder
}

func (builder *Builder[IT, OT]) Then(parser *Parser) *Builder {
	builder.parser = Then(builder.parser, parser)
	return builder
}

func (builder *Builder[IT, OT]) Skip(parser *Parser) *Builder {
	builder.parser = Skip(builder.parser, parser)
	return builder
}

func (builder *Builder[IT, OT]) AtMost(times int) *Builder {
	builder.parser = AtMost(builder.parser, times)
	return builder
}

func (builder *Builder[IT, OT]) AtLeast(times int) *Builder {
	builder.parser = AtLeast(builder.parser, times)
	return builder
}

func (builder *Builder[IT, OT]) Many() *Builder {
	builder.parser = Many(builder.parser)
	return builder
}

func (builder *Builder[IT, OT]) SeparatedBy(sepParser *Parser, min int, max int) *Builder {
	builder.parser = SeparatedBy(builder.parser, sepParser, min, max)
	return builder
}
