package apc

import "fmt"

type Origin struct {
	Name    string
	LineNum int
	ColNum  int
}

func (origin Origin) String() string {
	return fmt.Sprintf("%v:%v:%v", origin.Name, origin.LineNum, origin.ColNum)
}
