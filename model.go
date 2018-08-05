package main

import (
	"fmt"

	v "github.com/joushou/gocnc/vector"
)

// Line is a straight path from Start to End
type Line struct {
	From v.Vector
	To   v.Vector
}

// Move returns the start and end points of the line
func (l Line) Move() (v.Vector, v.Vector) {
	return l.From, l.To
}

// Reverse reverses the line, taking the bulge into account
func (l *Line) Reverse() {
	l.From, l.To = l.To, l.From
}

func (l Line) Equal(m Move) bool {
	if m, ok := m.(*Line); ok {
		return l.From == m.From && l.To == m.To
	}
	return false
}

func (l Line) String() string {
	return fmt.Sprintf("%s -> %s", l.From, l.To)
}

// Arc is an arc from Start to End around Center, either CW or CCW
type Arc struct {
	From   v.Vector
	To     v.Vector
	Center v.Vector
	CW     bool
}

// Move returns the start and end point of the arc
func (a *Arc) Move() (v.Vector, v.Vector) {
	return a.From, a.To
}

// Reverse reverses the line, taking the bulge into account
func (a *Arc) Reverse() {
	a.From, a.To = a.To, a.From
	a.CW = !a.CW
}

func (a Arc) Equal(m Move) bool {
	if m, ok := m.(*Arc); ok {
		return a.From == m.From &&
			a.To == m.To &&
			a.Center == m.Center &&
			a.CW == m.CW
	}
	return false
}

func (a *Arc) String() string {
	return fmt.Sprintf("%s ~> %s", a.From, a.To)
}
