package main

import (
	"fmt"
	"math"
)

const PRECISION = 5

type Vec struct {
	X, Y float64
}

func (v Vec) String() string {
	s := fmt.Sprintf("<%.2f,%.2f>", v.X, v.Y)
	return s
}

// addition of two vectors
func (v Vec) Add(w Vec) Vec {
	return Vec{v.X + w.X, v.Y + w.Y}
}

// substraction of two vectors
func (v Vec) Sub(w Vec) Vec {
	return Vec{v.X - w.X, v.Y - w.Y}
}

// norm of a vector
func (v Vec) Norm() float64 {
	return math.Sqrt(math.Pow(v.X, 2) + math.Pow(v.Y, 2))
}

// distance between two points
func (v Vec) Distance(w Vec) float64 {
	return v.Sub(w).Norm()
}

// tests equality with a given precision
func (v Vec) Equals(w Vec) bool {
	if v == w {
		return true
	}
	pre := math.Pow10(-PRECISION)
	dx := math.Abs(v.X - w.X)
	dy := math.Abs(v.Y - w.Y)
	return dx < pre && dy < pre
}

// angle of a vector in radians
func (v Vec) Angle() float64 {
	angle := math.Atan2(v.Y, v.X)
	if angle < 0 {
		angle += math.Pi * 2
	}
	return angle
}

// converts degrees to radians
func Radians(angle float64) float64 {
	return angle / 180 * math.Pi
}

// cartesian coordinates from polar
func Cartesian(angle float64, radius float64) Vec {
	p := Vec{}
	p.X = radius * math.Cos(angle)
	p.Y = radius * math.Sin(angle)
	return p
}

// polar coordinates (angle, radius) from cartesian
func Polar(p Vec) (float64, float64) {
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
func ArcToBulge(center Vec, radius float64, startAngle float64, endAngle float64) (Vec, Vec, float64) {
	startPoint := Cartesian(startAngle, radius).Add(center)
	endPoint := Cartesian(endAngle, radius).Add(center)
	// bulge conversion
	theta := endAngle - startAngle
	bulge := math.Tan(theta / 4)
	return startPoint, endPoint, bulge
}

// from a bulge representation of an arc, computes center, start angle, end angle and radius
// usefull for gcode generation
func BulgeToArc(p1 Vec, p2 Vec, bulge float64) (Vec, float64, float64, float64) {
	theta2 := 2 * math.Atan(bulge)                // half of included angle
	d := p1.Distance(p2) / 2                      // half of the chord length
	r := d / math.Sin(theta2)                     // radius
	a := p2.Sub(p1).Angle()                       // angle of the chord
	c := Cartesian(math.Pi/2-theta2+a, r).Add(p1) // center
	if bulge < 0 {
		return c, p2.Sub(c).Angle(), p1.Sub(c).Angle(), math.Abs(r)
	} else {
		return c, p1.Sub(c).Angle(), p2.Sub(c).Angle(), math.Abs(r)
	}
}
