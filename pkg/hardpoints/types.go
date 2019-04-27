package hardpoints

import (
	"fmt"
)

type Point struct {
	Path string
	Line int
	Tags []string
}

func (p Point) String() string {
	return fmt.Sprintf("%s:%d", p.Path, p.Line)
}
