package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

// WIP to try out other algorithm

type V struct{ X, Y float64 }

type NURBS struct {
	Degree  int
	Knots   []float64
	Weights []float64
	Ctrls   []V
}

func (n NURBS) FindSpan(u float64) int {
	idx := len(n.Knots) - n.Degree - 1
	// special case
	if u == n.Knots[idx+1] {
		return idx
	}
	// binary search
	U := n.Knots
	low := n.Degree
	high := idx + 1
	mid := (low + high) / 2
	for u < U[mid] || u >= U[mid+1] {
		if u < U[mid] {
			high = mid
		} else {
			low = mid
		}
		mid = (low + high) / 2
	}
	return mid
}

// i: index of the span
// p: degree of the B-Spline
// u: position
func (n NURBS) Basis(i, p int, u float64) []float64 {
	U := n.Knots
	N := make([]float64, p+1, p+1)
	L := make([]float64, p+1, p+1)
	R := make([]float64, p+1, p+1)

	N[0] = 1
	for j := 1; j <= p; j++ {
		L[j] = u - U[i+1-j] // distance to left bound
		R[j] = U[i+j] - u   // distance to right bound
		saved := 0.0
		for r := 0; r < j; r++ {
			tmp := N[r] / (R[r+1] + L[j-r])
			N[r] = saved + R[r+1]*tmp
			saved = L[j-r] * tmp
		}
		N[j] = saved
	}
	return N
}

func (n NURBS) Eval(u float64) V {
	return V{}
}

func (n NURBS) KnotBounds() (float64, float64) {
	return n.Knots[0], n.Knots[len(n.Knots)-1]
}

var nurbs = NURBS{
	Degree:  2,
	Knots:   []float64{0, 0, 0, 0.5, 1, 1, 1},
	Ctrls:   []V{V{0, 0}, V{-50, 112.5}, V{150, 112.5}, V{100, 0}},
	Weights: []float64{1, 1, 1, 1},
}

func main() {
	spew.Dump(nurbs)
	min, max := nurbs.KnotBounds()
	for i := 0; i < 11; i++ {
		u := min + float64(i)/10*(max-min)
		k := nurbs.FindSpan(u)
		fmt.Println("u", u)
		fmt.Println("span", k)
		fmt.Println(u, nurbs.Basis(k, nurbs.Degree, u))
	}
}
