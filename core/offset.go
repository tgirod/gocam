package core

import (
	"fmt"
	"math"
)

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
	var lines []*Line
	for _, l := range p.Lines {
		lines = append(lines, l.Offset(d))
	}

	offset := NewPath(fmt.Sprintf("%s off(%f)", p.Name, d))

	// next manage the transitions between lines
	for cur := 0; cur < p.Len(); cur++ {
		next := (cur + 1) % p.Len()
		// add the current
		offset.Append(lines[cur])
		// if offset lines are not connected
		if !lines[cur].End.Equals(lines[next].Start) {
			angle := p.Lines[next].TanStart().Angle() - p.Lines[cur].TanEnd().Angle()
			if angle > 0 {
				// add an arc to join lines
				center := p.Lines[cur].End
				radius := math.Abs(d)
				startAngle := lines[cur].End.Sub(center).Angle()
				endAngle := lines[next].Start.Sub(center).Angle()
				startPoint, endPoint, bulge := ArcToBulge(center, radius, startAngle, endAngle)
				arc := Line{
					Start: startPoint,
					End:   endPoint,
					Bulge: bulge}
				offset.Append(arc)
			} else {
				// cut lines to match at intersection
			}
		}
	}
	return offset
}
