package main

import (
	"github.com/joushou/gocnc/gcode"
)

// Path is a sequence of connected moves (Moves[i].To == Moves[i+1].From)
type Path []Move

func NewPath(m Move) Path {
	return Path{m}
}

// Move represents anything that moves from point A to point B and can be reversed
type Move interface {
	Move() (Vector, Vector)
	Reverse()
	Equal(Move) bool
}

// Move returns the start and end points of the path
func (p Path) Move() (Vector, Vector) {
	if len(p) == 0 {
		return Vector{}, Vector{}
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

func (p Path) Points() []Vector {
	pts := make([]Vector, 0, len(p)+1)
	for i, m := range p {
		from, to := m.Move()
		if i == 0 {
			pts = append(pts, from)
		}
		pts = append(pts, to)
	}
	return pts
}

const EPSILON float64 = 1E-3

func (p *Path) Append(m Move) bool {
	// utility function
	app := func() {
		if mPath, ok := m.(Path); ok {
			*p = append(*p, mPath...)
		} else {
			*p = append(*p, m)
		}
	}

	prep := func() {
		if mPath, ok := m.(Path); ok {
			*p = append(mPath, (*p)...)
		} else {
			*p = append(Path{m}, (*p)...)
		}
	}

	if p.IsClosed() {
		return false
	}

	// empty path, always append
	if len(*p) == 0 {
		app()
		return true
	}

	pFrom, pTo := p.Move()
	mFrom, mTo := m.Move()
	// append
	if pTo == mFrom {
		app()
		return true
	}
	// reverse append
	if pTo == mTo {
		m.Reverse()
		app()
		return true
	}
	// prepend
	if mTo == pFrom {
		prep()
		return true
	}
	// reverse prepend
	if mFrom == pFrom {
		m.Reverse()
		prep()
		return true
	}

	// no match, discard and return false
	return false
}

// IsClosed returns true if the path ends where it started
func (p *Path) IsClosed() bool {
	from, to := p.Move()
	return len(*p) > 0 && from == to
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

func (p Path) Gcode() []gcode.Block {
	bs := []gcode.Block{}

	// initial G0 move to the starting point
	start, _ := p.Move()
	bs = append(bs, move(start))

	// subsequent G1 moves
	for _, m := range p {
		// add actual move (G1, G2 or G3)
		if g, ok := m.(Gcoder); ok {
			bs = append(bs, g.Gcode())
		} else {
			Log.Printf("Move of type %T does not implement Gcoder", m)
		}
	}
	return bs
}
