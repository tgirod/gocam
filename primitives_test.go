package main

import (
	"testing"

	v "github.com/joushou/gocnc/vector"
	"github.com/stretchr/testify/assert"
)

func TestLineReverse(t *testing.T) {
	v, w := v.Vector{0, 0, 0}, v.Vector{1, 1, 0}
	a, b := &Line{v, w}, &Line{w, v}
	b.Reverse()

	assert.Equal(t, a, b, "Line.Reverse failed")
}

func TestLineEqual(t *testing.T) {
	l1 := &Line{v.Vector{0, 0, 0}, v.Vector{1, 1, 1}}
	l2 := &Line{v.Vector{0, 0, 0}, v.Vector{1, 1, 1}}
	assert.Equal(t, true, l1.Equal(l2), "should be equal")
}

func TestArcReverse(t *testing.T) {
	a := &Arc{v.Vector{1, 0, 0}, v.Vector{-1, 0, 0}, v.Vector{0, 0, 0}, false}
	b := &Arc{v.Vector{-1, 0, 0}, v.Vector{1, 0, 0}, v.Vector{0, 0, 0}, true}
	b.Reverse()
	assert.Equal(t, a, b, "Arc.Reverse failed")
}

func TestArcEqual(t *testing.T) {
	a1 := &Arc{v.Vector{-1, 0, 0}, v.Vector{1, 0, 0}, v.Vector{0, 0, 0}, false}
	a2 := &Arc{v.Vector{-1, 0, 0}, v.Vector{1, 0, 0}, v.Vector{0, 0, 0}, false}
	assert.Equal(t, true, a1.Equal(a2), "should be equal")
}
