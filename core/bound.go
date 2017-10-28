package core

import "math"

type Bounds struct {
	Lower Vector
	Upper Vector
}

// returns the bounding box surrounding the line.
// FIXME: this does not take the bulge into account
func (l *Line) Bounds() Bounds {
	lower := Vector{
		math.Min(l.From.X, l.To.X),
		math.Min(l.From.Y, l.To.Y),
	}
	upper := Vector{
		math.Max(l.From.X, l.To.X),
		math.Max(l.From.Y, l.To.Y),
	}
	return Bounds{lower, upper}
}

// returns the bounding box surrounding the path by combining the bounding
// boxes of each line
func (p *Path) Bounds() Bounds {
	bounds := Bounds{} // default value

	if p.Len() > 0 {
		bounds = p.Lines[0].Bounds()
		for _, l := range p.Lines[1:] {
			b := l.Bounds()
			b.Lower.X = math.Min(b.Lower.X, bounds.Lower.X)
			b.Lower.Y = math.Min(b.Lower.Y, bounds.Lower.Y)
			b.Upper.X = math.Max(b.Upper.X, bounds.Upper.X)
			b.Upper.Y = math.Max(b.Upper.Y, bounds.Upper.Y)
		}
	}
	return bounds
}
