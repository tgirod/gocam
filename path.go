package main

import v "github.com/joushou/gocnc/vector"

// Path is a sequence of connected moves (Moves[i].To == Moves[i+1].From)
type Path struct {
	Moves []Move
}

// Move represents anything that moves from point A to point B and can be reversed
type Move interface {
	Move() (v.Vector, v.Vector)
	Adjust(v.Vector, v.Vector)
	Reverse()
}

// Len returns the number of moves composing the path
func (p *Path) Len() int {
	return len(p.Moves)
}

// Move returns the start and end points of the path
func (p *Path) Move() (v.Vector, v.Vector) {
	if p.Len() == 0 {
		return v.Vector{}, v.Vector{}
	} else {
		from, _ := p.Moves[0].Move()
		_, to := p.Moves[p.Len()-1].Move()
		return from, to
	}
}

func (p *Path) Adjust(from, to v.Vector) {
	if p.Len() == 0 {
		return
	}
	// adjust the starting point
	first := p.Moves[0]
	_, toFirst := first.Move()
	first.Adjust(from, toFirst)
	// adjust the ending point
	last := p.Moves[p.Len()-1]
	fromLast, _ := last.Move()
	last.Adjust(fromLast, to)
}

// Reverse reverses path p, and all its composing moves
func (p *Path) Reverse() {
	i := 0
	j := p.Len() - 1

	for i <= j {
		if i == j {
			p.Moves[i].Reverse()
		} else {
			p.Moves[i].Reverse()
			p.Moves[j].Reverse()
			p.Moves[i], p.Moves[j] = p.Moves[j], p.Moves[i]
		}
		i++
		j--
	}
}

func (p *Path) Points() []v.Vector {
	pts := []v.Vector{}
	for i, m := range p.Moves {
		from, to := m.Move()
		if i == 0 {
			pts = append(pts, from)
		}
		pts = append(pts, to)
	}
	return pts
}

const EPSILON float64 = 0.001

// Append attempts to add a move at the end of the path. It will work if:
// * the path is empty
// * m.From == p.To
// * the distance between the two points is sufficiently small
// in those cases the move will be added and true is returned
// otherwise the path isn't changed and false is returned
func (p *Path) Append(m Move) bool {
	// empty path, always append
	if p.Len() == 0 {
		p.Moves = append(p.Moves, m)
		return true
	}
	_, pTo := p.Move()
	mFrom, mTo := m.Move()
	// exact match, append
	if pTo == mFrom {
		p.Moves = append(p.Moves, m)
		return true
	}
	// approximate match, append
	dist := pTo.Diff(mFrom).Norm()
	if dist < EPSILON {
		// adjust m.From to have an exact match
		m.Adjust(pTo, mTo)
		p.Moves = append(p.Moves, m)
		return true
	}
	// doesn't match, return false
	return false
}

// Join attempts to concatenate two paths
func (p *Path) Join(p2 *Path) bool {
	// empty path, always append
	if p.Len() == 0 {
		p.Moves = append(p.Moves, p2.Moves...)
		return true
	}
	_, pTo := p.Move()
	p2From, p2To := p2.Move()
	// exact match, append
	if pTo == p2From {
		p.Moves = append(p.Moves, p2.Moves...)
		return true
	}
	// approximate match, append
	dist := pTo.Diff(p2From).Norm()
	if dist < EPSILON {
		p2.Adjust(pTo, p2To)
		p.Moves = append(p.Moves, p2.Moves...)
		return true
	}
	// doesn't match, return false
	return false
}

// IsClosed returns true if the path ends where it started
func (p *Path) IsClosed() bool {
	from, to := p.Move()
	return from == to || from.Diff(to).Norm() < EPSILON
}

// IsClockwise returns true if the path is running clockwise, false otherwise.
// The shoelace algorithm is used to determine the direction of rotation
func (p *Path) IsClockwise() bool {
	sum := 0.0
	for _, m := range p.Moves {
		from, to := m.Move()
		sum += (to.X - from.X) * (to.Y + from.Y)
	}
	// the curve is CW if the sum is positive, CCW if the sum is negative
	return sum > 0
}

func (p *Path) Flatten() {
	flat := Path{}
	for _, m := range p.Moves {
		if m, ok := m.(*Path); ok {
			m.Flatten() // recursively flatten
			if ok := flat.Join(m); !ok {
				panic("this should not happen")
			}
		} else {
			if ok := flat.Append(m); !ok {
				panic("this should not happen")
			}
		}
	}
}
