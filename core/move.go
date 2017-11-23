package core

import (
	"fmt"
)

// Mover interface should be implemented by any type representing a move
type Mover interface {
	StartPoint() Vector
	EndPoint() Vector
	Reverse()
}

// Line is a straight path from Start to End
type Line struct {
	Start Vector
	End   Vector
}

// StartPoint returns the starting point
func (l *Line) StartPoint() Vector {
	return l.Start
}

// EndPoint returns the ending point
func (l *Line) EndPoint() Vector {
	return l.End
}

// Reverse reverses the line, taking the bulge into account
func (l *Line) Reverse() {
	l.Start, l.End = l.End, l.Start
}

// Vector returns the vector from Start to End
func (l *Line) Vector() Vector {
	return l.End.Sub(l.Start)
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %v, %v", l.Start, l.End)
}

// Arc is an arc from Start to End around Center, either CW or CCW
type Arc struct {
	Start  Vector
	End    Vector
	Center Vector
	CW     bool
}

// StartPoint returns the starting point
func (a *Arc) StartPoint() Vector {
	return a.Start
}

// EndPoint returns the ending point
func (a *Arc) EndPoint() Vector {
	return a.End
}

// Reverse reverses the line, taking the bulge into account
func (a *Arc) Reverse() {
	a.Start, a.End = a.End, a.Start
	a.CW = !a.CW
}

// Radius returns the radius of the arc
func (a *Arc) Radius() float64 {
	return a.Start.Sub(a.Center).Length()
}
