package core

import "fmt"

// Offset method returns a new line at distance d of line l
func (l *Line) Offset(d float64) *Line {
	var nl Line
	if l.Bulge == 0 {
		// straight line
		off := l.TanStart().Normal().Mul(d)
		nl.Start = l.Start.Add(off)
		nl.End = l.End.Add(off)
		nl.Bulge = 0
	} else if l.Bulge > 0 {
		// CCW arc
		offStart := l.TanStart().Normal().Mul(d)
		nl.Start = l.Start.Add(offStart)
		offEnd := l.TanEnd().Normal().Mul(d)
		nl.End = l.End.Add(offEnd)
		nl.Bulge = l.Bulge
	} else {
		// CW arc
		offStart := l.TanStart().Normal().Mul(-d)
		nl.Start = l.Start.Add(offStart)
		offEnd := l.TanEnd().Normal().Mul(-d)
		nl.End = l.End.Add(offEnd)
		nl.Bulge = l.Bulge
	}
	return &nl
}

// Offset method for a Path. returns a new path at distance d
func (p *Path) Offset(d float64) *Path {
	// first, generate an offset for each line
	offLines := NewPath(fmt.Sprintf("%s offset(%.2f)", p.Name, d))
	for _, l := range p.Lines {
		offLines.Append(l.Offset(d))
	}
	// next manage the transitions between lines
	// 1. if ends are already equals, leave it this way (it means the originals were tangent)
	// 2. there are two possibilities:
	// 	a. angle>180 degrees --> add an arc to join lines
	//  b. angle<180 degrees --> find the intersection and modify both lines to match there
	return offLines // FIXME
}
