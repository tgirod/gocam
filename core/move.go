package core

import (
	"fmt"
	"math"
)

// Line represents a line between two points From and To. The line can be
// either straight (if Bulge == 0) or an arc for others values of Bulge.
type Line struct {
	From  Vector
	To    Vector
	Bulge float64
}

// Reverse reverses the line, taking the bulge into account
func (l *Line) Reverse() {
	l.From, l.To = l.To, l.From
	l.Bulge = -l.Bulge
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %v, %v, %.2f", l.From, l.To, l.Bulge)
}

// EndAngle returns the angle of the tangent at the endpoint of line l
func (l *Line) EndAngle() float64 {
	var angle float64
	if l.Bulge == 0 {
		// straight line
		angle = l.To.Diff(l.From).Angle()
	} else if l.Bulge > 0 {
		// CCW arc
		_, _, angle, _ := BulgeToArc(l.From, l.To, l.Bulge)
		angle += math.Pi / 2
	} else {
		// CW arc
		_, _, angle, _ := BulgeToArc(l.From, l.To, l.Bulge)
		angle -= math.Pi / 2
	}
	return angle
}
