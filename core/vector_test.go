package core

import (
	"math"
	"testing"
)

var (
	a         = Vector{1, 1}
	b         = Vector{2, 2}
	f float64 = 2
)

func TestAdd(t *testing.T) {
	v := Vector{3, 3}
	if a.Add(b) != v {
		t.Error("Add failed")
	}
}

func TestSub(t *testing.T) {
	v := Vector{-1, -1}
	if a.Sub(b) != v {
		t.Error("Sub failed")
	}
}

func TestMul(t *testing.T) {
	v := Vector{2, 2}
	if a.Mul(f) != v {
		t.Error("Mul failed")
	}
}

func TestDiv(t *testing.T) {
	v := Vector{0.5, 0.5}
	if a.Div(f) != v {
		t.Error("Mul failed")
	}
}

func TestNeg(t *testing.T) {
	v := Vector{-1, -1}
	if a.Neg() != v {
		t.Error("Neg failed")
	}
}

func TestLength(t *testing.T) {
	res := math.Sqrt(2)
	if a.Length() != res {
		t.Error("Neg failed")
	}
}

func TestUnit(t *testing.T) {
	l := b.Unit().Length()
	if math.Abs(l-1) > tolerance {
		t.Error("Unit failed")
	}
}

func TestNormal(t *testing.T) {
	if !a.Normal().Equals(Vector{1, -1}) {
		t.Error("Normal failed")
	}
}
