package core

import (
	"fmt"
	"math"
)

const tolerance = 10E-5

// Vector represents a 2D vector
type Vector struct {
	X, Y float64
}

func (v Vector) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", v.X, v.Y)
}

// Add returns the addition v + o
func (v Vector) Add(o Vector) Vector {
	return Vector{
		X: v.X + o.X,
		Y: v.Y + o.Y}
}

// Sub returns the difference v - o
func (v Vector) Sub(o Vector) Vector {
	return Vector{
		X: v.X - o.X,
		Y: v.Y - o.Y}
}

// Mul multiplies the vector by scalar s
func (v Vector) Mul(s float64) Vector {
	return Vector{
		X: v.X * s,
		Y: v.Y * s}
}

// Div divides the vector by scalar s
func (v Vector) Div(s float64) Vector {
	return Vector{
		X: v.X / s,
		Y: v.Y / s}
}

// Neg returns the negative of vector v
func (v Vector) Neg() Vector {
	return Vector{
		X: -v.X,
		Y: -v.Y}
}

// Length returns the norm of v
func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normal returns the normal vector of v
func (v Vector) Normal() Vector {
	return Vector{
		X: v.Y,
		Y: -v.X}
}

// Unit returns vector v divided by v.Norm()
func (v Vector) Unit() Vector {
	n := v.Length()
	return v.Div(n)
}

// Equals tests if vectors v and o are sufficiently close to be considered
// equal. Returns true if the distance between the two vectors is below
// tolerance, otherwise false
func (v Vector) Equals(o Vector) bool {
	return v.Sub(o).Length() < tolerance
}

// Angle returns the angle (in radians) formed by the vector v
func (v Vector) Angle() float64 {
	angle := math.Atan2(v.Y, v.X)
	if angle < 0 {
		angle += math.Pi * 2
	}
	return angle
}

// Dot returns the dot product of vectors v and o
func (v Vector) Dot(o Vector) float64 {
	return v.X*o.X + v.Y*o.Y
}
