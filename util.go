package main

import (
	"log"
	"math"
	"os"
)

// Log is used to print logs to stderr.
var Log *log.Logger = log.New(os.Stderr, "Gocam: ", 0)

// Radians converts degrees to radians.
func deg2rad(angle float64) float64 {
	return angle / 180 * math.Pi
}

// Cartesian converts polar coordinates to cartesian.
func pol2car(angle float64, radius float64) Vector {
	p := Vector{}
	p.X = radius * math.Cos(angle)
	p.Y = radius * math.Sin(angle)
	return p
}

// Polar converts cartesian coordinates to polar (angle, radius).
func car2pol(p Vector) (float64, float64) {
	radius := math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))
	angle := math.Atan2(p.Y, p.X)
	return angle, radius
}

// ArcToBulge converts dxf arc representation (center, radius, and angles) to
// line+bulge representation (start point, end point, and bulge). Angles are in
// radians.
// Bulge == 0: straight line
// Bulge > 0: CCW arc
// Bulge < 0: CW arc
// Bulge == 1: semi-circle
func arcToBulge(center Vector, radius float64, startAngle float64, endAngle float64) (Vector, Vector, float64) {
	startPoint := pol2car(startAngle, radius).Sum(center)
	endPoint := pol2car(endAngle, radius).Sum(center)
	// bulge conversion
	theta := endAngle - startAngle
	bulge := math.Tan(theta / 4)
	return startPoint, endPoint, bulge
}

// BulgeToArc converts line+bulge arc representation to dxf representation.
func bulgeToArc(p1 Vector, p2 Vector, bulge float64) (Vector, float64, float64, float64) {
	theta2 := 2 * math.Atan(bulge) // half of included angle
	chord := p2.Diff(p1)           // chord of the arc
	d := chord.Norm() / 2          // half of the chord's length
	a := vec2angle(chord)          // angle of the chord
	r := d / math.Sin(theta2)      // radius

	center := pol2car(math.Pi/2-theta2+a, r).Sum(p1) // center
	radius := math.Abs(r)
	var startAngle, endAngle float64
	if bulge < 0 {
		startAngle = vec2angle(p2.Diff(center))
		endAngle = vec2angle(p1.Diff(center))
	} else {
		startAngle = vec2angle(p1.Diff(center))
		endAngle = vec2angle(p2.Diff(center))
	}
	return center, radius, startAngle, endAngle
}

func vec2angle(v Vector) float64 {
	return math.Atan2(v.Y, v.X)
}
