package core

// Offsetter defines the necessary interface to compute an offset
type Offsetter interface {
	Tangents() (Vector, Vector)
}

// Tangents returns the tangents of Line l at starting and ending points
func (l *Line) Tangents() (Vector, Vector) {
	t := l.Vector().Unit()
	return t, t
}

// TanStart returns the tangent vector at the starting point of arc a
func (a *Arc) Tangents() (Vector, Vector) {
	tanStart := a.Start.Sub(a.Center).Normal()
	tanEnd := a.End.Sub(a.Center).Normal()
	if !a.CW {
		tanStart = tanStart.Neg()
		tanEnd = tanEnd.Neg()
	}
	return tanStart, tanEnd
}

// Offset method returns a new line at distance d of line l
func (l *Line) Offset(d float64) *Line {
	t, _ := l.Tangents()
	off := t.Normal().Mul(d)
	return &Line{
		Start: l.Start.Add(off),
		End:   l.End.Add(off)}
}

// Offset method returns a new arc at distance d of arc a
func (a *Arc) Offset(d float64) *Arc {
	ts, te := a.Tangents()
	offStart := ts.Normal().Mul(d)
	offEnd := te.Normal().Mul(d)

	return &Arc{
		Start: a.Start.Add(offStart),
		End:   a.End.Add(offEnd),
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
