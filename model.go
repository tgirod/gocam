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
func (l *Line) Move() (v.Vector, v.Vector) {
	return l.From, l.To
}

// Reverse reverses the line, taking the bulge into account
func (l *Line) Reverse() {
	l.From, l.To = l.To, l.From
}

// Adjust replaces the start and end points of the line
func (l *Line) Adjust(from, to v.Vector) {
	l.From = from
	l.To = to
}

func (l *Line) String() string {
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

// Adjust replaces the start and end points of the arc
func (a *Arc) Adjust(from, to v.Vector) {
	a.From = from
	a.To = to
}

func (a *Arc) String() string {
	return fmt.Sprintf("%s ~> %s", a.From, a.To)
}
