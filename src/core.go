package main

import (
	"fmt"
	"strings"
)

type Mover interface {
	Start() Vec  // returns the coordinates of the move's start
	Reverse(Vec) // reverse the move toward the given coordinates
}

type Line struct {
	From  Vec
	Bulge float64
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %s,%.2f", l.From, l.Bulge)
}

func (l *Line) Start() Vec {
	return l.From
}

func (l *Line) Reverse(v Vec) {
	l.From = v
	l.Bulge = -l.Bulge
}

type Path struct {
	Handle string  // a name for this path
	Moves  []Mover // a sequence of moves to form a path
	End    Vec     // endpoint of the path
}

func NewPath(end Vec, handle string) *Path {
	return &Path{handle, []Mover{}, end}
}

func (p *Path) String() string {
	l := make([]string, p.Len()+2)
	l[0] = fmt.Sprintf("Path %s:", p.Handle)
	for i, m := range p.Moves {
		l[i+1] = fmt.Sprintf("\t%s", m)
	}
	l[p.Len()+1] = fmt.Sprintf("\tEnd: %s", p.End)
	return strings.Join(l, "\n")
}

func (p *Path) Len() int {
	return len(p.Moves)
}

func (p *Path) Start() Vec {
	if p.Len() > 0 {
		return p.Moves[0].Start()
	} else {
		return p.End
	}
}

func (p *Path) AppendMove(m Mover) {
	p.Moves = append(p.Moves, m)
}

func (p *Path) Join(q *Path) {
	p.Moves = append(p.Moves, q.Moves...)
	p.End = q.End
	p.Handle = fmt.Sprintf("%s->%s", p.Handle, q.Handle)
}

func (p *Path) Reverse() {
	if p.Len() > 0 {

		// reverse each move
		start := p.Start()
		for i := 0; i < p.Len()-1; i++ {
			p.Moves[i].Reverse(p.Moves[i+1].Start())
		}
		p.Moves[p.Len()-1].Reverse(p.End)
		p.End = start

		// reverse the order of the moves
		for i, j := 0, p.Len()-1; i < j; i, j = i+1, j-1 {
			p.Moves[i], p.Moves[j] = p.Moves[j], p.Moves[i]
		}
	}
}

func (p *Path) IsClosed() bool {
	return p.Start().Equals(p.End)
}

func (p *Path) IsClockwise() bool {
	sum := 0.0
	for i := 0; i < p.Len(); i++ {
		cur := p.Moves[i].Start()
		next := p.End
		if i < p.Len()-1 {
			next = p.Moves[i+1].Start()
		}
		prod := (next.X - cur.X) * (next.Y + cur.Y)
		sum += prod
	}
	// the curve is CW if the sum is positive, CCW if the sum is negative
	return sum > 0
}

type Document struct {
	Paths []Path
}

func (doc *Document) Len() int {
	return len(doc.Paths)
}

func (doc *Document) String() string {
	l := make([]string, doc.Len())
	for i, m := range doc.Paths {
		l[i] = fmt.Sprintf("%s", m)
	}
	return strings.Join(l, "\n")
}

// Join a path to an existing one if possible, otherwise append it to the list
func (doc *Document) AddPath(add *Path) {
	inserted := false
	for i := 0; i < doc.Len() && !inserted; i++ {
		cur := &doc.Paths[i]
		if !cur.IsClosed() {
			if cur.End.Equals(add.Start()) {
				// CUR -> ADD
				cur.Join(add)
				inserted = true
			} else if add.End.Equals(cur.Start()) {
				// ADD -> CUR
				add.Join(cur)
				doc.Paths[i] = *add
				inserted = true
			}
		}
	}

	if !inserted {
		// add cannot be joined to any existing path
		doc.Paths = append(doc.Paths, *add)
	}
}

func (doc *Document) RemovePath(idx int) {
	if idx < len(doc.Paths) {
		doc.Paths = append(doc.Paths[:idx], doc.Paths[idx+1:]...)
	}
}

func (doc *Document) Regroup() {
	for i := 0; i < len(doc.Paths); i++ {
		cur := &doc.Paths[i]
		if !cur.IsClosed() {
			end := cur.End
			// search for a path that starts at "end" and join it to "cur"
			for j, _ := range doc.Paths {
				path := &doc.Paths[j]
				start := path.Start()
				if !path.IsClosed() && start.Equals(end) {
					// joining paths
					cur.Join(path)
					doc.RemovePath(j)
					// a path has been joined to cur. Let's try again
					i--
					break
				}
			}
		}
	}
	for i, _ := range doc.Paths {
		doc.Paths[i].Handle = fmt.Sprintf("%d", i)
	}
}
