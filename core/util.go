package core

import (
	"log"
	"math"
	"os"
)

var Log *log.Logger = log.New(os.Stderr, "", 0)

// converts degrees to radians
func Radians(angle float64) float64 {
	return angle / 180 * math.Pi
}

// cartesian coordinates from polar
func Cartesian(angle float64, radius float64) Vector {
	p := Vector{}
	p.X = radius * math.Cos(angle)
	p.Y = radius * math.Sin(angle)
	return p
}

// polar coordinates (angle, radius) from cartesian
func Polar(p Vector) (float64, float64) {
	radius := math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))
	angle := math.Atan2(p.Y, p.X)
	return angle, radius
}

// convert classical arc representation to (start,end,bulge) format. Angles
// must be passed in radians.
// Bulge == 0: straight line
// Bulge > 0: CCW arc
// Bulge < 0: CW arc
// Bulge == 1: semi-circle
func ArcToBulge(center Vector, radius float64, startAngle float64, endAngle float64) (Vector, Vector, float64) {
	startPoint := Cartesian(startAngle, radius).Sum(center)
	endPoint := Cartesian(endAngle, radius).Sum(center)
	// bulge conversion
	theta := endAngle - startAngle
	bulge := math.Tan(theta / 4)
	return startPoint, endPoint, bulge
}

// from a bulge representation of an arc, computes center, start angle, end angle and radius
// usefull for gcode generation
func BulgeToArc(p1 Vector, p2 Vector, bulge float64) (Vector, float64, float64, float64) {
	theta2 := 2 * math.Atan(bulge)                // half of included angle
	d := p2.Diff(p1).Norm() / 2                   // half of the chord length
	r := d / math.Sin(theta2)                     // radius
	a := p2.Diff(p1).Angle()                      // angle of the chord
	c := Cartesian(math.Pi/2-theta2+a, r).Sum(p1) // center
	if bulge < 0 {
		return c, p2.Diff(c).Angle(), p1.Diff(c).Angle(), math.Abs(r)
	} else {
		return c, p1.Diff(c).Angle(), p2.Diff(c).Angle(), math.Abs(r)
	}
}