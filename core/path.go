package core

import (
	"fmt"
	"strings"
)

type Path struct {
	Name  string // a name for this path
	Lines []Line // a sequence of lines to form a path
}

func NewPath(name string) *Path {
	return &Path{name, []Line{}}
}

func (p *Path) String() string {
	s := make([]string, p.Len()+1)
	s[0] = fmt.Sprintf("Path %s:", p.Name)
	for i, l := range p.Lines {
		s[i+1] = fmt.Sprintf("\t%s", l)
	}
	return strings.Join(s, "\n")
}

func (p *Path) Len() int {
	return len(p.Lines)
}

func (p *Path) Start() Vector {
	return p.Lines[0].Start()
}

func (p *Path) End() Vector {
	return p.Lines[p.Len()-1].End()
}

func (p *Path) Append(l Line) {
	p.Lines = append(p.Lines, l)
}

func (p *Path) Join(q *Path) {
	p.Lines = append(p.Lines, q.Lines...)
	p.Name = fmt.Sprintf("%s->%s", p.Name, q.Name)
}

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

func (p *Path) IsClosed() bool {
	return p.Start().Equals(p.End())
}

// This method is using the shoelace algorithm to determine if the path is
// clockwise or counter-clockwise.
func (p *Path) IsClockwise() bool {
	sum := 0.0
	for _, l := range p.Lines {
		start := l.Start()
		end := l.End()
		sum += (end.X - start.X) * (end.Y + start.Y)
	}
	// the curve is CW if the sum is positive, CCW if the sum is negative
	return sum > 0
}

// If the path is closed, StartNear(v) will choose the vertex closest to v as a
// starting point for the path
func (p *Path) StartNear(v Vector) {
	if p.Len() > 0 && p.IsClosed() {
		// find the closest vertex
		index := 0
		closest := v.Diff(p.Lines[0].Start()).Norm()
		for i := 1; i < p.Len(); i++ {
			current := v.Diff(p.Lines[i].Start()).Norm()
			if current < closest {
				index = i
				closest = current
			}
		}
		// rotate
		p.Lines = append(p.Lines[index:], p.Lines[:index]...)
	}
}
