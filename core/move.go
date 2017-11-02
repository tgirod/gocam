package core

import (
	"fmt"
)

// Line represents a line between two points From and To. The line can be
// either straight (if Bulge == 0) or an arc for others values of Bulge.
type Line struct {
	Start Vector
	End   Vector
	Bulge float64
}

// Reverse reverses the line, taking the bulge into account
func (l *Line) Reverse() {
	l.Start, l.End = l.End, l.Start
	l.Bulge = -l.Bulge
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %v, %v, %.2f", l.Start, l.End, l.Bulge)
}
