package main

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
		term1 := 0.0
		if t[i+k] > t[i] {
			term1 = (u - t[i]) / (t[i+k] - t[i]) * s.n(i, k-1, u)
		}
		term2 := 0.0
		if t[i+k+1] > t[i+1] {
			term2 = (t[i+k+1] - u) / (t[i+k+1] - t[i+1]) * s.n(i+1, k-1, u)
		}
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
		v := P[i].Multiply(fact)
		res = res.Sum(v)
		div += fact
	}
	return res.Divide(div)
}

/*

This is not used yet, as I'm still unsure about how it works. I think the
principle is this:

0. evaluate NURBS at position u
1. find k, the index of the knot span associated to u (findKnotSpan)
2. get the respective influence of the local control points (findKnotSpan)
3. get the local control points (how ?)
4. compute the weighted mean between the control points
*/

// based on https://github.com/mfem/mfem/blob/master/mesh/nurbs.cpp#L214
// func (s Spline) findKnotSpan(u float64) int {
// 	order := s.Degree
// 	var low, mid, high int
// 	if u == s.Knots[len(s.Controls)+order] {
// 		mid = len(s.Controls)
// 	} else {
// 		low = order
// 		high = len(s.Controls) + 1
// 		mid = (low + high) / 2
// 		for (u < s.Knots[mid-1]) || (u > s.Knots[mid]) {
// 			if u < s.Knots[mid-1] {
// 				high = mid
// 			} else {
// 				low = mid
// 			}
// 			mid = (low + high) / 2
// 		}
// 	}
// 	return mid
// }

// based on https://www.researchgate.net/publication/228411721/
func (s Spline) basisITS0(k, p int, u float64) []float64 {
	N := make([]float64, p, p)     // FIXME
	L := make([]float64, p+1, p+1) // FIXME
	R := make([]float64, p+1, p+1) // FIXME

	N[0] = 1
	for j := 1; j <= p; j++ {
		saved := 0.0
		L[j] = u - s.Knots[k+1-j] // distance to left bound
		R[j] = s.Knots[k+j] - u   // distance to right bound
		for r := 0; r < j; r++ {
			tmp := N[r] / (R[r+1] + L[j-r])
			N[r] = saved + R[r+1]*tmp
			saved = L[j-r] * tmp
		}
		N[j] = saved
	}
	return N
}

func (s Spline) eval2(u float64) Vector {
	P := s.Controls
	n := len(s.Controls)
	w := s.Weights
	k := s.Degree
	res := Vector{}
	div := 0.0
	for i := 0; i < n; i++ {
		fact := w[i] * s.n(i, k, u)
		v := P[i].Multiply(fact)
		res = res.Sum(v)
		div += fact
	}
	return res.Divide(div)
}

func (s Spline) greville() []float64 {
	order := s.Degree + 1
	greLength := len(s.Knots) - order + 1
	gre := make([]float64, greLength, greLength)
	for i := 0; i < len(gre); i++ {
		for j := i; j < i+order; j++ {
			gre[i] += s.Knots[j]
		}
		gre[i] = gre[i] / float64(order)
	}
	return gre
}
