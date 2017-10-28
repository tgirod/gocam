package core

import (
	"fmt"
	"math"
)

const TOLERANCE = 10E-5

type Vector struct {
	X, Y float64
}

func (v Vector) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", v.X, v.Y)
}

func (v Vector) Sum(o Vector) Vector {
	return Vector{
		X: v.X + o.X,
		Y: v.Y + o.Y}
}

func (v Vector) Diff(o Vector) Vector {
	return Vector{
		X: v.X - o.X,
		Y: v.Y - o.Y}
}

func (v Vector) Norm() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func (v Vector) Equals(o Vector) bool {
	return v.Diff(o).Norm() < TOLERANCE
}

func (v Vector) Angle() float64 {
	angle := math.Atan2(v.Y, v.X)
	if angle < 0 {
		angle += math.Pi * 2
	}
	return angle
}
