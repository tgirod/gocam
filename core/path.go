package core

import (
	"fmt"
	"strings"
)

// Path is a continuous sequence of lines, ie Lines[i].To == Lines[i+1].From
type Path struct {
	Name  string // a name for this path
	Lines []Line // a sequence of lines to form a path
}

// NewPath creates a new empty path
func NewPath(name string) *Path {
	return &Path{name, []Line{}}
}

func (p *Path) String() string {
	s := make([]string, p.Len()+1)
	s[0] = fmt.Sprintf("Path %s:", p.Name)
	for i, l := range p.Lines {
		s[i+1] = fmt.Sprintf("\t%v", l)
	}
	return strings.Join(s, "\n")
}

// Len returns the number of lines in the path
func (p *Path) Len() int {
	return len(p.Lines)
}

// Start returns the starting point of a non-empty path
func (p *Path) Start() Vector {
	return p.Lines[0].Start
}

// End returns the ending point of a non-empty path
func (p *Path) End() Vector {
	return p.Lines[p.Len()-1].End
}

// Append adds a line at the end of the path
func (p *Path) Append(l *Line) {
	p.Lines = append(p.Lines, *l)
}

// Join appends the path q to the path p
func (p *Path) Join(q *Path) {
	p.Lines = append(p.Lines, q.Lines...)
	p.Name = fmt.Sprintf("%s->%s", p.Name, q.Name)
}

// Reverse reverts path p by calling reverse on each line, then reversing the
// order of the lines
func (p *Path) Reverse() {
	if p.Len() > 0 {
		// reversing each line
		for i := 0; i < p.Len(); i++ {
			p.Lines[i].Reverse()
		}

		// reversing the order of the lines
		for i, j := 0, p.Len()-1; i < j; i, j = i+1, j-1 {
			p.Lines[i], p.Lines[j] = p.Lines[j], p.Lines[i]
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
	return p.Start().Equals(p.End())
}

// IsClockwise returns true if the path is running clockwise, false otherwise.
// The shoelace algorithm is used to determine the direction of rotation
func (p *Path) IsClockwise() bool {
	sum := 0.0
	for _, l := range p.Lines {
		sum += (l.End.X - l.Start.X) * (l.End.Y + l.Start.Y)
	}
	// the curve is CW if the sum is positive, CCW if the sum is negative
	return sum > 0
}

// StartNear modifies the order of the lines of path p so the starting point is
// as close as possible to vector v. If the path is open, nothing is changed.
func (p *Path) StartNear(v Vector) {
	if p.Len() > 0 && p.IsClosed() {
		// find the closest vertex
		index := 0
		closest := v.Sub(p.Lines[0].Start).Length()
		for i := 1; i < p.Len(); i++ {
			current := v.Sub(p.Lines[i].Start).Length()
			if current < closest {
				index = i
				closest = current
			}
		}
		// rotate
		p.Lines = append(p.Lines[index:], p.Lines[:index]...)
	}
}
