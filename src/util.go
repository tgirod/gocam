package main

import (
	"math"

	"github.com/rpaloschi/dxf-go/core"
)

const PRECISION = 5

type Point struct {
	core.Point
}

func (p Point) Equals(q Point) bool {
	if p == q {
		return true
	}
	pre := math.Pow10(-PRECISION)
	dx := math.Abs(p.X - q.X)
	dy := math.Abs(p.Y - q.Y)
	dz := math.Abs(p.Z - q.Z)
	return dx < pre && dy < pre && dz < pre
}

func pol2car(center Point, radius float64, angle float64) Point {
	p := Point{}
	p.X = center.X + radius*math.Cos(angle/180*math.Pi)
	p.Y = center.Y + radius*math.Sin(angle/180*math.Pi)
	p.Z = center.Z
	return p
}
