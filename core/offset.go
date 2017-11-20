package core

// Offsetter defines the necessary interface to compute an offset
type Offsetter interface {
	TanStart() Vector
	TanEnd() Vector
}

// TanStart returns the tangent vector at the starting point of line l
func (l *Line) TanStart() Vector {
	return l.Vector().Unit()
}

// TanEnd returns the tangent vector at the ending point of line l
func (l *Line) TanEnd() Vector {
	return l.Vector().Unit()
}

// TanStart returns the tangent vector at the starting point of arc a
func (a *Arc) TanStart() Vector {
	vec := a.Start.Sub(a.Center).Normal()
	if !a.CW {
		vec = vec.Neg()
	}
	return vec
}

// TanEnd returns the tangent vector at the ending point of arc a
func (a *Arc) TanEnd() Vector {
	vec := a.End.Sub(a.Center).Normal()
	if !a.CW {
		vec = vec.Neg()
	}
	return vec
}

// Offset method returns a new line at distance d of line l
func (l *Line) Offset(d float64) *Line {
	var nl Line
	// straight line
	off := l.TanStart().Normal().Mul(d)
	nl.Start = l.Start.Add(off)
	nl.End = l.End.Add(off)
	return &nl
}

// Offset method returns a new arc at distance d of arc a
func (a *Arc) Offset(d float64) *Arc {
	radius := a.Radius()
	ratio := (radius + d) / radius

	return &Arc{
		Start: a.Start.Sub(a.Center).Mul(ratio).Add(a.Center),
		End:   a.End.Sub(a.Center).Mul(ratio).Add(a.Center),
		CW:    a.CW}
}

// Offset method for a Path. returns a new path at distance d
// func (p *Path) Offset(d float64) *Path {
// 	var off []Offset
// 	for _, m := range p.Moves {
// 		if o, ok := m.(Offset); ok {
// 			off = append(off, o.Offset(d))
// 		}
// 	}

// 	offset := NewPath(fmt.Sprintf("%s off(%f)", p.Name, d))

// 	// next manage the transitions between lines
// 	for cur := 0; cur < p.Len(); cur++ {
// 		next := (cur + 1) % p.Len()
// 		// add the current
// 		offset.Append(lines[cur])
// 		// if offset lines are not connected
// 		if !lines[cur].End.Equals(lines[next].Start) {
// 			angle := p.Lines[next].TanStart().Angle() - p.Lines[cur].TanEnd().Angle()
// 			if angle > 0 {
// 				// add an arc to join lines
// 				center := p.Lines[cur].End
// 				radius := math.Abs(d)
// 				startAngle := lines[cur].End.Sub(center).Angle()
// 				endAngle := lines[next].Start.Sub(center).Angle()
// 				startPoint, endPoint, bulge := ArcToBulge(center, radius, startAngle, endAngle)
// 				arc := Line{
// 					Start: startPoint,
// 					End:   endPoint,
// 					Bulge: bulge}
// 				offset.Append(arc)
// 			} else {
// 				// cut lines to match at intersection
// 			}
// 		}
// 	}
// 	return offset
// }
