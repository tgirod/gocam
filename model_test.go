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

func TestLineAdjust(t *testing.T) {
	l := Line{
		From: v.Vector{1, 1, 0},
		To:   v.Vector{2, 2, 0},
	}
	z := v.Vector{0, 0, 0}
	l.Adjust(z, z)
	from, to := l.Move()
	assert.Equal(t, from, z, "Line.Adjust failed")
	assert.Equal(t, to, z, "Line.Adjust failed")
}

func TestArcReverse(t *testing.T) {
	a := &Arc{v.Vector{1, 0, 0}, v.Vector{-1, 0, 0}, v.Vector{0, 0, 0}, false}
	b := &Arc{v.Vector{-1, 0, 0}, v.Vector{1, 0, 0}, v.Vector{0, 0, 0}, true}
	b.Reverse()
	assert.Equal(t, a, b, "Arc.Reverse failed")
}

func TestArcAdjust(t *testing.T) {
	a := Arc{
		From:   v.Vector{1, 0, 0},
		To:     v.Vector{-1, 0, 0},
		Center: v.Vector{0, 0, 0},
		CW:     false,
	}

	z := v.Vector{0, 0, 0}
	a.Adjust(z, z)
	from, to := a.Move()
	assert.Equal(t, from, z, "Arc.Adjust failed")
	assert.Equal(t, to, z, "Arc.Adjust failed")
}
