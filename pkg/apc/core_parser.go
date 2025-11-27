package apc

import "fmt"

func Range[IT, OT any](parser Parser[IT, OT], min int, max int) Parser[IT, []OT] {
	if min < 0 {
		panic("cannot use Range parser with min < 0")
	}
	if max < 0 && max != -1 {
		panic("cannot use Range parser with max < 0 and max != -1 (infinite)")
	}
	if min > max && max != -1 {
		panic("cannot use Range parser with min > max")
	}

	return NewCallbackParser(func(ctx Context[IT]) ([]OT, error) {
		var err error
		results := []OT{}
		for len(results) < max || max == -1 {
			var result OT
			result, err = parser.Parse(ctx)
			if err != nil {
				break
			}
			results = append(results, result)
		}
		if len(results) < min {
			if err != nil {
				return nil, err
			}
			return nil, ExpErr{
				AtIndex:    ctx.CurIndex(),
				Expected:   fmt.Sprintf("%d to %d of %s", min, max, parser),
				Unexpected: "...unknown...",
			}
		}
		return results, nil
	})
}

func Seq[IT, OT any](parsers ...Parser[IT, OT]) Parser[IT, []OT] {
	if len(parsers) < 2 {
		panic("cannot have Seq Parser with < 2 parsers")
	}
	return NewCallbackParser(func(ctx Context[IT]) ([]OT, error) {
		results := []OT{}
		for _, parser := range parsers {
			result, err := parser.Parse(ctx)
			if err != nil {
				return nil, err
			}
			results = append(results, result)
		}
		return results, nil
	})
}

func Index[IT, OT any](index int, parser Parser[IT, []OT]) Parser[IT, OT] {
	if index < 0 {
		panic("the Index Parser 'index' must be >= 0")
	}
	return NewCallbackParser(func(ctx Context[IT]) (OT, error) {
		results, err := parser.Parse(ctx)
		if err != nil {
			return newT[OT](), err
		}
		if index >= len(results) {
			panic("invalid Index Parser 'index' was too large for sub-parser results")
		}
		return results[index], nil
	})
}

func Map[IT, OTA, OTB any](parser Parser[IT, OTA], transform func(value OTA) (OTB, error)) Parser[IT, OTB] {
	return NewCallbackParser(func(ctx Context[IT]) (OTB, error) {
		result, err := parser.Parse(ctx)
		if err != nil {
			return newT[OTB](), err
		}
		return transform(result)
	})
}
