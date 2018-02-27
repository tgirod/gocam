package model

import (
	"math"
	"testing"
)

var (
	v1, v2, v3 Vector = Vector{1, 1}, Vector{2, 2}, Vector{3, 3}
)

func TestAdd(t *testing.T) {
	v := Vector{3, 3}
	if v1.Add(v2) != v {
		t.Error("Add failed")
	}
}

func TestSub(t *testing.T) {
	v := Vector{-1, -1}
	if v1.Sub(v2) != v {
		t.Error("Sub failed")
	}
}

func TestMul(t *testing.T) {
	v := Vector{2, 2}
	if v1.Mul(2.) != v {
		t.Error("Mul failed")
	}
}

func TestDiv(t *testing.T) {
	v := Vector{0.5, 0.5}
	if v1.Div(2) != v {
		t.Error("Mul failed")
	}
}

func TestNeg(t *testing.T) {
	v := Vector{-1, -1}
	if v1.Neg() != v {
		t.Error("Neg failed")
	}
}

func TestLength(t *testing.T) {
	res := math.Sqrt(2)
	if v1.Length() != res {
		t.Error("Neg failed")
	}
}
