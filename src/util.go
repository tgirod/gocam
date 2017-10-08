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

func (v Vec) Add(w Vec) Vec {
	return Vec{v.X + w.X, v.Y + w.Y}
}

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

func (v Vec) Equals(w Vec) bool {
	if v == w {
		return true
	}
	pre := math.Pow10(-PRECISION)
	dx := math.Abs(v.X - w.X)
	dy := math.Abs(v.Y - w.Y)
	return dx < pre && dy < pre
}

// converts degrees to radians
func deg2rad(angle float64) float64 {
	return math.Cos(angle / 180 * math.Pi)
}

// convert polar to cartesian coordinates
func pol2car(angle float64, radius float64) Vec {
	p := Vec{}
	p.X = radius * math.Cos(angle)
	p.Y = radius * math.Sin(angle)
	return p
}

// convert cartesian to polar coordinates (angle, radius)
func car2pol(p Vec) (float64, float64) {
	radius := math.Sqrt(math.Pow(p.X, 2) + math.Pow(p.Y, 2))
	angle := math.Atan2(p.Y, p.X)
	return angle, radius
}

// computes the angle of the vector formed by those two points
func angle(a Vec, b Vec) float64 {
	v := b.Sub(a)
	return math.Atan2(v.Y, v.X)
}

// from a bulge representation of an arc, computes center, start angle, end angle and radius
// usefull for gcode generation
func convertBulge(p1 Vec, p2 Vec, bulge float64) (Vec, float64, float64, float64) {
	theta2 := 2 * math.Atan(bulge)              // half of included angle
	d := p1.Distance(p2) / 2                    // half of the chord length
	r := d / math.Sin(theta2)                   // radius
	a := angle(p1, p2)                          // angle of the chord
	c := pol2car(math.Pi/2-theta2+a, r).Add(p1) // center
	if bulge < 0 {
		return c, angle(c, p2), angle(c, p1), math.Abs(r)
	} else {
		return c, angle(c, p1), angle(c, p2), math.Abs(r)
	}
}
