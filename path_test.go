package main

import (
	"testing"

	v "github.com/joushou/gocnc/vector"
	"github.com/stretchr/testify/assert"
)

var a, b, c, d = v.Vector{0, 0, 0}, v.Vector{1, 0, 0}, v.Vector{2, 0, 0}, v.Vector{3, 0, 0}
var e = v.Vector{1, 1, 0}
var ab, bc, cd = &Line{a, b}, &Line{b, c}, &Line{c, d}
var ba, cb, dc = &Line{b, a}, &Line{c, b}, &Line{d, c}
var c2 = c.Sum(v.Vector{EPSILON / 2, 0, 0})

func path(points ...v.Vector) *Path {
	p := Path{}
	for i := 0; i < len(points)-1; i++ {
		from := points[i]
		to := points[i+1]
		line := &Line{from, to}
		p = append(p, line)
	}
	return &p
}

func TestLen0(t *testing.T) {
	p := Path{}
	assert.Equal(t, 0, len(p), "length of empty path is not zero")
}

func TestEqual(t *testing.T) {
	p1 := path(a, b, c, d)
	p2 := path(a, b, c, d)
	assert.Equal(t, true, p1.Equal(p2), "should be equal")
}

func TestPoints(t *testing.T) {
	p := Path{
		ab,
		bc,
	}
	pts := []v.Vector{a, b, c}
	assert.Equal(t, pts, p.Points(), "Path.Points failed")
}

func TestAppendEmpty(t *testing.T) {
	p := &Path{}
	ok := p.Append(ab)
	assert.Equal(t, ok, true, "Path.Append failed to append to an empty path")
	assert.Equal(t, &Path{ab}, p, "Path.Append didn't append the right thing")
}

func TestAppendExact(t *testing.T) {
	p := &Path{ab}
	ok := p.Append(bc)
	assert.Equal(t, ok, true, "cannot append move with shared end")
}

func TestAppendApproximate(t *testing.T) {
	p := &Path{bc}
	ok := p.Append(&Line{c2, d})
	assert.Equal(t, true, ok, "cannot append move with approximate end")
}

func TestAppendFail(t *testing.T) {
	p := &Path{ab}
	ok := p.Append(cd)
	assert.Equal(t, false, ok, "should not append")
	assert.Equal(t, &Path{ab}, p, "should be ab")
}

func TestReverseEmpty(t *testing.T) {
	p := &Path{}
	p.Reverse()
	assert.Equal(t, true, p.Equal(&Path{}), "empty reverse")
}

func TestReverseEven(t *testing.T) {
	p := path(a, b, c)
	p.Reverse()
	assert.Equal(t, path(c, b, a), p, "Path.Reverse failed")
}

func TestReverseOdd(t *testing.T) {
	p := path(a, b, c, d)
	p.Reverse()
	assert.Equal(t, path(d, c, b, a), p, "Path.Reverse failed")
}

func TestIsClosed(t *testing.T) {
	p := path(a, b, c, a)
	ok := p.IsClosed()
	assert.Equal(t, true, ok, "Path.IsClosed failed")
}

func TestIsClockwise(t *testing.T) {
	p := path(a, b, e, a)
	ok := p.IsClockwise()
	assert.Equal(t, false, ok, "Path.IsClockwise failed to detect CCW path")
}

func TestIsCounterClockwise(t *testing.T) {
	p := path(a, e, b, a)
	ok := p.IsClockwise()
	assert.Equal(t, true, ok, "Path.IsClockwise failed to detect CCW path")
}

// func TestHasInnerLoop(t *testing.T) {
// 	var data = []struct {
// 		in        *Path
// 		hasLoop   bool
// 		startLoop int
// 		endLoop   int
// 		doc       string
// 	}{
// 		{poly(a, b, c, d), false, 0, 0, "no loop"},
// 		{poly(a, b, c, d, a), false, 0, 0, "outer loop"},
// 		{poly(a, b, c, a, d), true, 0, 3, "start loop"},
// 		{poly(a, b, c, d, b), true, 1, 4, "end loop"},
// 		{poly(a, b, c, b, d), true, 1, 3, "inner loop"},
// 	}
// 	for _, d := range data {
// 		hasLoop, startLoop, endLoop := d.in.HasInnerLoop()
// 		if d.hasLoop != hasLoop || d.startLoop != startLoop || d.endLoop != endLoop {
// 			t.Error("Path.HasInnerLoop failed", d.doc, hasLoop, startLoop, endLoop)
// 		}
// 	}
// }

// func TestSplitInnerLoops(t *testing.T) {
// 	var data = []struct {
// 		in  *Path
// 		out []*Path
// 		doc string
// 	}{
// 		{poly(a, b, c, d), []*Path{poly(a, b, c, d)}, "no loop"},
// 		{poly(a, b, c, d, a), []*Path{poly(a, b, c, d, a)}, "outer loop"},
// 		{poly(a, b, c, a, d), []*Path{poly(a, b, c, a), poly(a, d)}, "start loop"},
// 		{poly(a, b, c, d, b), []*Path{poly(b, c, d, b), poly(a, b)}, "end loop"},
// 		{poly(a, b, c, b, d), []*Path{poly(b, c, b), poly(a, b, d)}, "inner loop"},
// 	}
// 	for _, d := range data {
// 		paths := d.in.SplitInnerLoops()
// 		if len(paths) != len(d.out) {
// 			t.Error("Path.SplitInnerLoops failed: wrong number of loops")
// 		}
// 		for i := 0; i < len(paths); i++ {
// 			if !paths[i].Equals(d.out[i]) {
// 				t.Error("Path.SplitInnerLoops failed on", d.doc, "test", i)
// 				t.Log("RETURNED", paths[i])
// 				t.Log("EXPECTED", d.out[i])
// 			}
// 		}
// 	}
// }
