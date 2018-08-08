package main

import (
	"fmt"

	"github.com/joushou/gocnc/gcode"
)

// Line is a straight path from Start to End
type Line struct {
	From Vector
	To   Vector
}

// Move returns the start and end points of the line
func (l Line) Move() (Vector, Vector) {
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
	return fmt.Sprintf("Line<(%.2f,%.2f)--(%.2f,%.2f)>", l.From.X, l.From.Y, l.To.X, l.To.Y)
}

func (l Line) Gcode() gcode.Block {
	g := &gcode.Block{}
	g.AppendNode(word('G', 1))
	g.AppendNodes(xy(l.To)...)
	return *g
}

// Arc is an arc from Start to End around Center, either CW or CCW
type Arc struct {
	From   Vector
	To     Vector
	Center Vector
	CW     bool
}

// Move returns the start and end point of the arc
func (a Arc) Move() (Vector, Vector) {
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

func (a Arc) String() string {
	return fmt.Sprintf("Arc<(%.2f,%.2f)--(%.2f, %.2f)>", a.From.X, a.From.Y, a.To.X, a.To.Y)
}

func (a Arc) Gcode() gcode.Block {
	b := &gcode.Block{}
	if a.CW {
		b.AppendNode(word('G', 2))
	} else {
		b.AppendNode(word('G', 3))
	}

	// arc's endpoint
	b.AppendNodes(xy(a.To)...)
	// center (relative to the start)
	center := a.Center.Diff(a.From)
	b.AppendNodes(ij(center)...)
	return *b
}

type Spline struct {
	Degree   int
	Closed   bool
	Knots    []float64
	Controls []Vector
	Weights  []float64
}

func (s Spline) Move() (Vector, Vector) {
	// FIXME
	return Vector{}, Vector{}
}

func (s *Spline) Reverse() {
	// FIXME
}

func (s Spline) Equal(m Move) bool {
	// FIXME
	if _, ok := m.(*Spline); ok {
		return true
	}
	return false
}
