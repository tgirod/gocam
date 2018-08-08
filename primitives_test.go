package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineReverse(t *testing.T) {
	v, w := Vector{0, 0, 0}, Vector{1, 1, 0}
	a, b := &Line{v, w}, &Line{w, v}
	b.Reverse()

	assert.Equal(t, a, b, "Line.Reverse failed")
}

func TestLineEqual(t *testing.T) {
	l1 := &Line{Vector{0, 0, 0}, Vector{1, 1, 1}}
	l2 := &Line{Vector{0, 0, 0}, Vector{1, 1, 1}}
	assert.Equal(t, true, l1.Equal(l2), "should be equal")
}

func TestArcReverse(t *testing.T) {
	a := &Arc{Vector{1, 0, 0}, Vector{-1, 0, 0}, Vector{0, 0, 0}, false}
	b := &Arc{Vector{-1, 0, 0}, Vector{1, 0, 0}, Vector{0, 0, 0}, true}
	b.Reverse()
	assert.Equal(t, a, b, "Arc.Reverse failed")
}

func TestArcEqual(t *testing.T) {
	a1 := &Arc{Vector{-1, 0, 0}, Vector{1, 0, 0}, Vector{0, 0, 0}, false}
	a2 := &Arc{Vector{-1, 0, 0}, Vector{1, 0, 0}, Vector{0, 0, 0}, false}
	assert.Equal(t, true, a1.Equal(a2), "should be equal")
}
