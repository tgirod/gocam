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

// Vector returns the vector from Start to End
func (l *Line) Vector() Vector {
	return l.End.Diff(l.Start)
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %v, %v, %.2f", l.Start, l.End, l.Bulge)
}

// TanStart returns the tangent vector at the starting point of the line
func (l *Line) TanStart() Vector {
	// simple case of the straight line
	if l.Bulge == 0 {
		return l.End.Diff(l.Start).Unit()
	}
	// FIXME arc
	return Vector{}
}

// TanEnd returns the tangent vector at the ending point of the line
func (l *Line) TanEnd() Vector {
	// simple case of the straight line
	if l.Bulge == 0 {
		return l.End.Diff(l.Start).Unit()
	}
	// FIXME arc
	return Vector{}
}
