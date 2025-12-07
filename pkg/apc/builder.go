package apc

type Builder struct {
	parser *Parser
}

func NewBuilder(initialParser *Parser) *Builder {
	return &Builder{
		parser: initialParser,
	}
}

func (builder *Builder) Build() *Parser {
	return builder.parser
}

func (builder *Builder) Describe(description string) *Builder {
	builder.parser = Describe(builder.parser, description)
	return builder
}

func (builder *Builder) Peek() *Builder {
	builder.parser = Peek(builder.parser)
	return builder
}

func (builder *Builder) Map(transform MapFunc) *Builder {
	builder.parser = Map(builder.parser, transform)
	return builder
}

func (builder *Builder) MapValue(value any) *Builder {
	builder.parser = MapValue(builder.parser, value)
	return builder
}

func (builder *Builder) MapSlice(transform MapSliceFunc) *Builder {
	builder.parser = MapSlice(builder.parser, transform)
	return builder
}

func (builder *Builder) Generate(generator ParserGeneratorFunc) *Builder {
	builder.parser = Generate(builder.parser, generator)
	return builder
}

func (builder *Builder) Index(index int) *Builder {
	builder.parser = Index(builder.parser, index)
	return builder
}

func (builder *Builder) ConcatSlices(parser *Parser) *Builder {
	builder.parser = ConcatSlices(builder.parser, parser)
	return builder
}

func (builder *Builder) Times(min int, max int) *Builder {
	builder.parser = Times(builder.parser, min, max)
	return builder
}

func (builder *Builder) Optional(defaultValue any) *Builder {
	builder.parser = Optional(builder.parser, defaultValue)
	return builder
}

func (builder *Builder) Head() *Builder {
	builder.parser = Head(builder.parser)
	return builder
}

func (builder *Builder) Tail() *Builder {
	builder.parser = Tail(builder.parser)
	return builder
}

func (builder *Builder) Then(parser *Parser) *Builder {
	builder.parser = Then(builder.parser, parser)
	return builder
}

func (builder *Builder) Skip(parser *Parser) *Builder {
	builder.parser = Skip(builder.parser, parser)
	return builder
}

func (builder *Builder) AtMost(times int) *Builder {
	builder.parser = AtMost(builder.parser, times)
	return builder
}

func (builder *Builder) AtLeast(times int) *Builder {
	builder.parser = AtLeast(builder.parser, times)
	return builder
}

func (builder *Builder) Many() *Builder {
	builder.parser = Many(builder.parser)
	return builder
}

func (builder *Builder) SeparatedBy(sepParser *Parser, min int, max int) *Builder {
	builder.parser = SeparatedBy(builder.parser, sepParser, min, max)
	return builder
}
