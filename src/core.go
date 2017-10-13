package main

import (
	"fmt"
	"strings"
)

type Mover interface {
	End() Vec    // returns the coordinates of the move's endpoint
	Reverse(Vec) // reverse the move toward the given coordinates
}

type Line struct {
	To    Vec
	Bulge float64
}

func (l *Line) String() string {
	return fmt.Sprintf("Line: %s,%.2f", l.To, l.Bulge)
}

func (l *Line) End() Vec {
	return l.To
}

func (l *Line) Reverse(v Vec) {
	l.To = v
	l.Bulge = -l.Bulge
}

type Path struct {
	Handle string  // a name for this path
	Moves  []Mover // a sequence of moves to form a path
	Start  Vec     // startpoint of the path
}

func NewPath(start Vec, handle string) *Path {
	return &Path{handle, []Mover{}, start}
}

func (p *Path) String() string {
	l := make([]string, p.Len()+2)
	l[0] = fmt.Sprintf("Path %s:", p.Handle)
	l[1] = fmt.Sprintf("\tStart: %s", p.Start)
	for i, m := range p.Moves {
		l[i+2] = fmt.Sprintf("\t%s", m)
	}
	return strings.Join(l, "\n")
}

func (p *Path) Len() int {
	return len(p.Moves)
}

func (p *Path) End() Vec {
	if p.Len() > 0 {
		return p.Moves[p.Len()-1].End()
	} else {
		return p.Start
	}
}

func (p *Path) AppendMove(m Mover) {
	p.Moves = append(p.Moves, m)
}

func (p *Path) Join(q *Path) {
	p.Moves = append(p.Moves, q.Moves...)
	p.Handle = fmt.Sprintf("%s->%s", p.Handle, q.Handle)
}

func (p *Path) Reverse() {
	if p.Len() > 0 {
		// saving endpoint before overwriting
		start := p.Moves[p.Len()-1].End()

		// reversing each move
		p.Moves[0].Reverse(p.Start)
		for i := 1; i < p.Len(); i++ {
			p.Moves[i].Reverse(p.Moves[i-1].End())
		}
		p.Start = start

		// reversing the order of the moves
		for i, j := 0, p.Len()-1; i < j; i, j = i+1, j-1 {
			p.Moves[i], p.Moves[j] = p.Moves[j], p.Moves[i]
		}

		// reversing Handle
		h := strings.Split(p.Handle, "->")
		for i, j := 0, len(h)-1; i < j; i, j = i+1, j-1 {
			h[i], h[j] = h[j], h[i]
		}
		p.Handle = strings.Join(h, "->")
	}
}

func (p *Path) IsClosed() bool {
	return p.Start.Equals(p.End())
}

// This method is using the shoelace algorithm to determine if the path is
// clockwise or counter-clockwise.
func (p *Path) IsClockwise() bool {
	sum := 0.0
	var cur, prev Vec
	for i := 0; i < p.Len(); i++ {
		if i == 0 {
			prev = p.Start
		} else {
			prev = p.Moves[i-1].End()
		}
		cur = p.Moves[i].End()

		prod := (cur.X - prev.X) * (cur.Y + prev.Y)
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
		l[i] = fmt.Sprintf("%s", &m)
	}
	return strings.Join(l, "\n")
}

func (doc *Document) AppendPath(p *Path) {
	doc.Paths = append(doc.Paths, *p)
}

func (doc *Document) RemovePath(idx int) {
	doc.Paths = append(doc.Paths[:idx], doc.Paths[idx+1:]...)
}

func (doc *Document) ExtractPath(idx int) *Path {
	p := &doc.Paths[idx]
	doc.Paths = append(doc.Paths[:idx], doc.Paths[idx+1:]...)
	return p
}

// Add a path to the document. This method will try to join the new path to the
// existing ones by looking for possible prepending and appending paths. Paths
// will be joined as necessary.
func (doc *Document) JoinPath(path *Path) {
	if path.IsClosed() {
		doc.AppendPath(path)
	} else {
		// search for prepending and appending paths
		pre := -1
		post := -1

		for i := 0; i < doc.Len() && (pre == -1 || post == -1); i++ {
			cur := &doc.Paths[i]
			if !cur.IsClosed() && cur.End().Equals(path.Start) {
				// found prepending path
				pre = i
			}
			if !cur.IsClosed() && path.End().Equals(cur.Start) {
				// found appending path
				post = i
			}
		}

		if pre == -1 && post == -1 {
			// isolated path
			// xxx -> PATH -> xxx
			doc.AppendPath(path)
		} else if post == -1 {
			// prepending path only
			// PRE -> PATH
			doc.Paths[pre].Join(path)
		} else if pre == -1 {
			// appending path only
			// PATH -> POST
			path.Join(&doc.Paths[post])
			doc.Paths[post] = *path
		} else if pre == post {
			// loop-closing path
			// PREPOST -> PATH
			doc.Paths[pre].Join(path)
		} else {
			// prepending and appending paths
			// PRE -> PATH -> POST
			doc.Paths[pre].Join(path)
			doc.Paths[pre].Join(&doc.Paths[post])
			doc.RemovePath(post)
		}
	}
}
