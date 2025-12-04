# APC - Another Parser Combinator

APC (Another Parser Combinator) is a minimalist parser combinator library written in Go.

APC used to rely on Go generics, but in the V3 redesign, generics are no longer used. This makes the flow and Go syntax of creating a parser easier to read and use overall.

APC is flexible enough to parse an input stream directly, or optionally the input can first be passed through an APC parser acting as a lexer. The tokens produced by this lexer can then be used as an input stream to token-based parsers.

APC supports backtracking by use of the `Peek` parser.

## Embedded Golang Struct BNF-Style Parse Tree Generator (`apcgen`)

The `apcgen` package provides support for generating parse trees by embedding the BNF-style syntax to define a parse tree into Golang structs. This can be done by directly matching on text input, or by providing a custom token-based parser.

This specific feature is inspired by [alecthomas/participle](https://github.com/alecthomas/participle), however the implementation is completely unique to `apc`. For more information on how to use `apcgen`, see [here](#apcgen-usage).
