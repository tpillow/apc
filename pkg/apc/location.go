package apc

import "fmt"

type Location struct {
	Name    string
	Index   int
	LineNum int
	ColNum  int
}

func (loc Location) String() string {
	return fmt.Sprintf("%s %d:%d:(#%d)", loc.Name, loc.LineNum, loc.ColNum, loc.Index)
}

func (loc Location) next(streamItem any) Location {
	lineNum := loc.LineNum
	colNum := loc.ColNum
	if r, ok := streamItem.(rune); ok {
		if r == '\n' {
			lineNum += 1
			colNum = 0
		} else {
			colNum += 1
		}
	}

	return Location{
		Name:    loc.Name,
		Index:   loc.Index + 1,
		LineNum: lineNum,
		ColNum:  colNum,
	}
}
