package core

import (
	"fmt"
	"math"
)

type Mover interface {
	Start() Vector // returns the coordinates of the move's startpoint
	End() Vector   // returns the coordinates of the move's endpoint
	Reverse()      // reverse the move
}

type Line struct {
	From  Vector
	To    Vector
	Bulge float64
}

func (l *Line) Start() Vector {
	return l.From
}

func (l *Line) End() Vector {
	return l.To
}

func (l *Line) Reverse() {
	l.From, l.To = l.To, l.From
	l.Bulge = -l.Bulge
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %v, %v, %.2f", l.From, l.To, l.Bulge)
}

// this method returns the angle formed by the tangent at the endpoint of the
// line.
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
