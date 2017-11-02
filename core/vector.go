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

// Sum returns the addition v + o
func (v Vector) Sum(o Vector) Vector {
	return Vector{
		X: v.X + o.X,
		Y: v.Y + o.Y}
}

// Diff returns the difference v - o
func (v Vector) Diff(o Vector) Vector {
	return Vector{
		X: v.X - o.X,
		Y: v.Y - o.Y}
}

// Norm returns the norm of v
func (v Vector) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Equals tests if vectors v and o are sufficiently close to be considered
// equal.
// returns true if the distance between the two vectors is below tolerance,
// otherwise false
func (v Vector) Equals(o Vector) bool {
	return v.Diff(o).Norm() < tolerance
}

// Angle returns the angle (in radians) formed by the vector v
func (v Vector) Angle() float64 {
	angle := math.Atan2(v.Y, v.X)
	if angle < 0 {
		angle += math.Pi * 2
	}
	return angle
}
