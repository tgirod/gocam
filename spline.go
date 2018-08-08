package main

import "fmt"

// De Boor's B-Spline evaluation algorithm, based on
// https://en.wikipedia.org/wiki/De_Boor's_algorithm#Example_implementation
// k: index of knot interval that contains x
// x: position
// t: array of knot positions, needs to be padded
// c: array of control points
// p: degree of B-spline
func deBoor(k int, x float64, t []float64, c []Vector, p int) Vector {
	// copy control points p+1 control points
	d := make([]Vector, p+1, p+1)
	copy(d, c[k-p:k+1])

	for r := 1; r < p+1; r++ {
		for j := p; j > r-1; j-- {
			alpha := (x - t[j+k-p]) / (t[j+1+k-r] - t[j+k-p])
			d[j] = d[j-1].Multiply(1 - alpha).Sum(d[j].Multiply(alpha))
		}
	}
	return d[p]
}

// implementation based on this ressource:
// http://web.cs.wpi.edu/~matt/courses/cs563/talks/nurbs.html

func (s Spline) n(i, k int, u float64) float64 {
	t := s.Knots
	if k == 0 {
		// special case to end recursion
		if t[i] <= u && u < t[i+1] {
			return 1
		} else {
			return 0
		}
	} else {
		// general case
		term1 := (u - t[i]) / (t[i+k] - t[i]) * s.n(i, k-1, u)
		term2 := (t[i+k+1] - u) / (t[i+k+1] - t[i+1]) * s.n(i+1, k-1, u)
		return term1 + term2
	}
}

func (s Spline) eval(u float64) Vector {
	P := s.Controls
	n := len(s.Controls)
	w := s.Weights
	k := s.Degree
	res := Vector{}
	div := 0.0
	for i := 0; i < n; i++ {
		fact := w[i] * s.n(i, k, u)
		fmt.Println(i, fact)
		res.Sum(P[i].Multiply(fact))
		div += fact
	}
	res.Divide(div)
	return res
}
