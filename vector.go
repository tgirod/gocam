package main

import "math"
import "fmt"

type Vector struct {
	X, Y float64
}

func (v Vector) Dot(o Vector) float64 {
	return v.X*o.X + v.Y*o.Y
}

func (v Vector) Norm() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector) Sum(o Vector) Vector {
	return Vector{
		X: v.X + o.X,
		Y: v.Y + o.Y,
	}
}

func (v Vector) Diff(o Vector) Vector {
	return Vector{
		X: v.X - o.X,
		Y: v.Y - o.Y,
	}
}

func (v Vector) Divide(d float64) Vector {
	return Vector{
		X: v.X / d,
		Y: v.Y / d,
	}
}

func (v Vector) Multiply(d float64) Vector {
	return Vector{
		X: v.X * d,
		Y: v.Y * d,
	}
}

func (v Vector) String() string {
	return fmt.Sprintf("(%f, %f)", v.X, v.Y)
}
