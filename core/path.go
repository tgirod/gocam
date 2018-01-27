package core

import (
	"fmt"
	"strings"
)

// Path is a continuous sequence of moves, ie Moves[i].End == Moves[i+1].Start
type Path struct {
	Name  string  // a name for this path
	Moves []Mover // a sequence of moves to form a path
}

// NewPath creates a new empty path
func NewPath(name string) *Path {
	return &Path{name, []Mover{}}
}

func (p *Path) String() string {
	s := make([]string, p.Len()+1)
	s[0] = fmt.Sprintf("Path %s:", p.Name)
	for i, m := range p.Moves {
		s[i+1] = fmt.Sprintf("\t%v", m)
	}
	return strings.Join(s, "\n")
}

// Len returns the number of moves in the path
func (p *Path) Len() int {
	return len(p.Moves)
}

// StartPoint returns the starting point of a non-empty path
func (p *Path) StartPoint() Vector {
	return p.Moves[0].StartPoint()
}

// EndPoint returns the ending point of a non-empty path
func (p *Path) EndPoint() Vector {
	return p.Moves[p.Len()-1].EndPoint()
}

// Append adds Move m at the end of Path p if their respective end and start
// are equals
func (p *Path) Append(m Mover) error {
	var err error
	if p.Len() > 0 && !p.EndPoint().Equals(m.StartPoint()) {
		err = fmt.Errorf("%v != %v", p.EndPoint(), m.StartPoint())
	} else {
		p.Moves = append(p.Moves, m)
	}
	return err
}

// Join appends the path q to the path p
func (p *Path) Join(q *Path) error {
	var err error
	if p.Len() > 0 && !p.EndPoint().Equals(q.StartPoint()) {
		err = fmt.Errorf("%v != %v", p.EndPoint(), q.StartPoint())
	} else {
		p.Moves = append(p.Moves, q.Moves...)
		p.Name = fmt.Sprintf("%s->%s", p.Name, q.Name)
	}
	return err
}

// Reverse reverts path p by calling reverse on each Move, then reversing the
// order of the moves
func (p *Path) Reverse() {
	if p.Len() > 0 {
		// reversing each Move
		for i := 0; i < p.Len(); i++ {
			p.Moves[i].Reverse()
		}

		// reversing the order of the Moves
		for i, j := 0, p.Len()-1; i < j; i, j = i+1, j-1 {
			p.Moves[i], p.Moves[j] = p.Moves[j], p.Moves[i]
		}

		// reversing name
		h := strings.Split(p.Name, "->")
		for i, j := 0, len(h)-1; i < j; i, j = i+1, j-1 {
			h[i], h[j] = h[j], h[i]
		}
		p.Name = strings.Join(h, "->")
	}
}

// IsClosed returns true if the starting and ending points of the path are
// equals
func (p *Path) IsClosed() bool {
	return p.StartPoint().Equals(p.EndPoint())
}

// IsClockwise returns true if the path is running clockwise, false otherwise.
// The shoelace algorithm is used to determine the direction of rotation
func (p *Path) IsClockwise() bool {
	sum := 0.0
	for _, m := range p.Moves {
		start := m.StartPoint()
		end := m.EndPoint()
		sum += (end.X - start.X) * (end.Y + start.Y)
	}
	// the curve is CW if the sum is positive, CCW if the sum is negative
	return sum > 0
}

// StartFrom modifies the order of the Moves of Path p so the starting point is
// as close as possible to vector v. If the path is open, nothing is changed.
func (p *Path) StartFrom(v Vector) {
	if p.Len() > 0 && p.IsClosed() {
		// find the closest vertex
		index := 0
		closest := v.Sub(p.Moves[0].StartPoint()).Length()
		for i := 1; i < p.Len(); i++ {
			dist := v.Sub(p.Moves[i].StartPoint()).Length()
			if dist < closest {
				index = i
				closest = dist
			}
		}
		// rotate
		p.Moves = append(p.Moves[index:], p.Moves[:index]...)
	}
}
