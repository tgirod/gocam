package model

import (
	"fmt"
)

// Move interface should be implemented by any type representing a move
type Move interface {
	Start() Vector
	End() Vector
	Reverse()
	Equals(m Move) bool
	Copy() Move
}

// Line is a straight path from Start to End
type Line struct {
	From Vector
	To   Vector
}

func (l *Line) Equals(m Move) bool {
	l2, ok := m.(*Line)
	return ok && *l == *l2
}

func (l *Line) Copy() Move {
	return &Line{l.From, l.To}
}

// Start returns the starting point
func (l *Line) Start() Vector {
	return l.From
}

// End returns the ending point
func (l *Line) End() Vector {
	return l.To
}

// Reverse reverses the line, taking the bulge into account
func (l *Line) Reverse() {
	l.From, l.To = l.To, l.From
}

// Vector returns the vector from Start to End
func (l *Line) Vector() Vector {
	return l.To.Sub(l.From)
}

func (l *Line) String() string {
	return fmt.Sprintf("<LINE: FROM=%v TO=%v>", l.From, l.To)
}

// Arc is an arc from Start to End around Center, either CW or CCW
type Arc struct {
	From   Vector
	To     Vector
	Center Vector
	CW     bool
}

func (a *Arc) Equals(m Move) bool {
	a2, ok := m.(*Arc)
	return ok && *a == *a2
}

func (a *Arc) Copy() Move {
	return &Arc{a.From, a.To, a.Center, a.CW}
}

// Start returns the starting point
func (a *Arc) Start() Vector {
	return a.From
}

// End returns the ending point
func (a *Arc) End() Vector {
	return a.To
}

// Reverse reverses the line, taking the bulge into account
func (a *Arc) Reverse() {
	a.From, a.To = a.To, a.From
	a.CW = !a.CW
}

func (a *Arc) Vector() Vector {
	return a.To.Sub(a.From)
}

func (a *Arc) String() string {
	return fmt.Sprintf("<ARC: FROM=%v TO=%v CENTER=%v CW=%v>", a.From, a.To, a.Center, a.CW)
}

// Radius returns the radius of the arc
func (a *Arc) Radius() float64 {
	return a.From.Sub(a.Center).Length()
}
