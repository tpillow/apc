# APC - Another Parser Combinator

APC (Another Parser Combinator) is a minimalist parser combinator library written in Go. APC utilizes Go generics as much as possible, to reduce the amount of type casting that must happen, as some other libraries rely on.

APC is flexible enough to parse an input stream directly, or optionally the input can first be passed through an APC parser acting as a lexer. The tokens produced by this lexer can then be used as an input stream to token-based parsers.

APC does not yet support backtracking/lookahead. My primary goal is to first get a stable combinator library setup before implementing lookahead support.

## Future Plans

- Better error messages / traces.
- Lookahead support.
- Write more tests.
- Doc generation (doc comments already added).

## Example Executables

- [JSON Parser](examples/json/main.go) that parses the JSON format into Go nodes.
- [JSON Parser using Lexer](examples/json_using_lexer/main.go) is the same as JSON Parser above, but it first lexes the input into tokens. This JSON parser then parses tokens instead of the raw input stream. This approach allows for better error messages.
- [Calculator](examples/calculator/main.go) that parses simple mathematical expression. With PEMDAS :)

## The Basics

### `Parser[CT, T]`

A parser is defined as a function with the following signature: `type Parser[CT, T any] func(ctx Context[CT]) (T, error)`. In other words, a function that takes a `Context[CT]` (where `CT` is the type of input stream token) for peeking/consuming an input stream, returning a tuple `(T, error)` where `T` is the return type of a successful parse. If parsing is successful, `error` should be `nil`. Otherwise, a parser may return one of 3 types of errors that must be handled:

1. `ParseError` - parsing failed, but input was NOT consumed. This means other parsers later in the line will be tried.
2. `ParseErrorConsumed` - parsing failed, and some input WAS consumed. This will cause an immediate fail of parsing entirely, as no lookahead is implemented yet.
3. Other error types - treated like `ParseErrorConsumed`. Could be I/O, buffer, etc. errors.

### Basic Parsers

TODO

### Combinator Parsers

TODO

### The `Ref` Parser

Creates a `Parser[CT, T]` from a `*Parser[CT, T]`. This is useful for avoiding circular dependencies. For example, the following is invalid due to a circular reference:

```go
// `value` refers to `hashValue`.
var value = Any[rune, any]("", MapToAny(ExactStr("hello")), MapToAny(hashValue))
// `hashValue` refers to `value`.
var hashValue = Seq[rune, any]("", MapToAny(ExactStr("#")), value)
```

However this can be remedied by using `Ref`:

```go
// `value` is just a variable declaration with no assignment.
var value Parser[rune, any]
// `valueRef` is a parser referring to `value`, which is not yet assigned.
var valueRef = Ref[rune, any](&value)
// `hashValue` refers to `valueRef` - NOT `value`.
var hashValue = Seq[rune, any]("", MapToAny(ExactStr("#")), valueRef)

func init() {
    // At runtime, `value` can then be defined and refer to `hashValue`:
    value = Any[rune, any]("", MapToAny(ExactStr("hello")), MapToAny(hashValue))
}
```

Note that in the above example `MapToAny(hashValue)` is necessary because a `Seq[CT, any]` returns `[]any` (not `any`).

### Origin

The `Origin` type holds information about a location in the input stream. This includes a name (usually the source filename), along with a line number and column number.

The current `Origin` of the input stream can be accessed by `Context.GetCurOrigin()`, and any type of `ParseError` will usually contain the `Origin` where the error originated.
