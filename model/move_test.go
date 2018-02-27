package model

import "testing"

func TestLineEqual(t *testing.T) {
	v, w := Vector{0, 0}, Vector{1, 1}
	a, b := &Line{v, w}, &Line{v, w}
	if !a.Equals(b) {
		t.Error("LineEqual failed")
		t.Log(a)
		t.Log(b)
		t.Log(a.Equals(b))
	}
}

func TestLineReverse(t *testing.T) {
	v, w := Vector{0, 0}, Vector{1, 1}
	a, b := &Line{v, w}, &Line{w, v}
	b.Reverse()

	if !a.Equals(b) {
		t.Error("LineReverse failed")
	}
}

func TestArcEqual(t *testing.T) {
	a := &Arc{Vector{1, 0}, Vector{-1, 0}, Vector{0, 0}, false}
	b := &Arc{Vector{1, 0}, Vector{-1, 0}, Vector{0, 0}, false}

	if !a.Equals(b) {
		t.Error("ArcEqual failed")
	}
}

func TestArcReverse(t *testing.T) {
	a := &Arc{Vector{1, 0}, Vector{-1, 0}, Vector{0, 0}, false}
	b := &Arc{Vector{-1, 0}, Vector{1, 0}, Vector{0, 0}, true}
	b.Reverse()

	if !a.Equals(b) {
		t.Error("ArcReverse failed")
	}
}
