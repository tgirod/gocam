package main

import (
	v "github.com/joushou/gocnc/vector"
)

// Path is a sequence of connected moves (Moves[i].To == Moves[i+1].From)
type Path []Move

// Move represents anything that moves from point A to point B and can be reversed
type Move interface {
	Move() (v.Vector, v.Vector)
	Reverse()
	Equal(Move) bool
}

// Move returns the start and end points of the path
func (p Path) Move() (v.Vector, v.Vector) {
	if len(p) == 0 {
		return v.Vector{}, v.Vector{}
	} else {
		from, _ := p[0].Move()
		_, to := p[len(p)-1].Move()
		return from, to
	}
}

// Reverse reverses path p, and all its composing moves
func (p Path) Reverse() {
	i := 0
	j := len(p) - 1

	for i <= j {
		if i == j {
			p[i].Reverse()
		} else {
			p[i].Reverse()
			p[j].Reverse()
			p[i], p[j] = p[j], p[i]
		}
		i++
		j--
	}
}

func (p Path) Equal(m Move) bool {
	if m, ok := m.(Path); ok {
		if len(p) != len(m) {
			return false
		}

		for i := range p {
			if !p[i].Equal(m[i]) {
				return false
			}
		}
	}

	return true
}

func (p Path) Points() []v.Vector {
	pts := make([]v.Vector, 0, len(p)+1)
	for i, m := range p {
		from, to := m.Move()
		if i == 0 {
			pts = append(pts, from)
		}
		pts = append(pts, to)
	}
	return pts
}

const EPSILON float64 = 0.001

func (p *Path) Append(m Move) bool {
	// utility function
	app := func() {
		if mPath, ok := m.(Path); ok {
			*p = append(*p, mPath...)
		} else {
			*p = append(*p, m)
		}
	}

	// empty path, always append
	if len(*p) == 0 {
		app()
		return true
	}

	_, pTo := p.Move()
	mFrom, _ := m.Move()
	// exact match, append
	if pTo == mFrom {
		app()
		return true
	}

	// approximate match, append
	dist := pTo.Diff(mFrom).Norm()
	if dist < EPSILON {
		app()
		return true
	}

	// no match, discard and return false
	return false
}

// IsClosed returns true if the path ends where it started
func (p *Path) IsClosed() bool {
	from, to := p.Move()
	return from == to || from.Diff(to).Norm() < EPSILON
}

// IsClockwise returns true if the path is running clockwise, false otherwise.
// The shoelace algorithm is used to determine the direction of rotation
func (p Path) IsClockwise() bool {
	sum := 0.0
	for _, m := range p {
		from, to := m.Move()
		sum += (to.X - from.X) * (to.Y + from.Y)
	}
	// the curve is CW if the sum is positive, CCW if the sum is negative
	return sum > 0
}
