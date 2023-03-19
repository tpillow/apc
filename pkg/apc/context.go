package apc

import "strings"

type RunePeeker interface {
	PeekRune(offset int) (rune, error)
}

type RuneConsumer interface {
	ConsumeRune() (rune, error)
}

type RunePeekerConsumer interface {
	RunePeeker
	RuneConsumer
}

type OriginGetter interface {
	GetOrigin() Origin
}

type Context interface {
	RunePeekerConsumer
	OriginGetter
}

// Return value may be less the `len` runes long.
func PeekNRunes(ctx Context, offset int, len int) (string, error) {
	sb := strings.Builder{}
	for i := offset; i < offset+len; i++ {
		r, err := ctx.PeekRune(i)
		if err != nil {
			return sb.String(), err
		}
		sb.WriteRune(r)
	}
	return sb.String(), nil
}

// Return value may be less the `len` runes long.
func ConsumeNRunes(ctx Context, len int) (string, error) {
	sb := strings.Builder{}
	for i := 0; i < len; i++ {
		r, err := ctx.ConsumeRune()
		if err != nil {
			return sb.String(), err
		}
		sb.WriteRune(r)
	}
	return sb.String(), nil
}
