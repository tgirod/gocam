package main

import (
	"testing"

	v "github.com/joushou/gocnc/vector"
	"github.com/stretchr/testify/assert"
)

var a, b, c, d = v.Vector{-1, -1, 0}, v.Vector{1, -1, 0}, v.Vector{1, 1, 0}, v.Vector{-1, 1, 0}
var c2 = v.Vector{1, 1 + EPSILON/2, 0}

var p Path = Path{}

func path(points ...v.Vector) Path {
	p := Path{}
	for i := 0; i < len(points)-1; i++ {
		from := points[i]
		to := points[i+1]
		line := &Line{from, to}
		p.Append(line)
	}
	return p
}

func TestLen0(t *testing.T) {
	assert.Equal(t, p.Len(), 0, "length of empty path is not zero")
}

func TestAppendEmpty(t *testing.T) {
	l := &Line{a, b}
	ok := p.Append(l)

	assert.Equal(t, ok, true, "cannot append move to empty path")
}

func TestAppendExact(t *testing.T) {
	l := &Line{b, c}
	ok := p.Append(l)

	assert.Equal(t, ok, true, "cannot append move with shared end")
}

func TestAppendApproximate(t *testing.T) {
	l := &Line{c2, d}
	ok := p.Append(l)

	assert.Equal(t, ok, true, "cannot append move with approximate end")
	lAdjusted := p.Moves[p.Len()-1]
	from, _ := lAdjusted.Move()
	assert.Equal(t, from, c, "From was not adjusted")
}

func TestPoints(t *testing.T) {
	pts := p.Points()
	exp := []v.Vector{a, b, c, d}

	assert.Equal(t, exp, pts, "unexpected points")
}

func TestReverseOdd(t *testing.T) {
	p := path(a, b, c)
	p.Reverse()
	revPts := p.Points()
	expPts := []v.Vector{c, b, a}

	assert.Equal(t, expPts, revPts, "Path.Reverse failed")
}

func TestReverseEven(t *testing.T) {
	p := path(a, b)
	p.Reverse()
	revPts := p.Points()
	expPts := []v.Vector{b, a}

	assert.Equal(t, expPts, revPts, "Path.Reverse failed")
}

func TestIsClosed(t *testing.T) {
	p := path(a, b, c, a)
	ok := p.IsClosed()
	assert.Equal(t, true, ok, "Path.IsClosed failed")
}

func TestIsClockwise(t *testing.T) {
	p := path(a, b, c, a)
	ok := p.IsClockwise()
	assert.Equal(t, false, ok, "Path.IsClockwise failed to detect CCW path")
}

func TestIsCounterClockwise(t *testing.T) {
	p := path(a, c, b, a)
	ok := p.IsClockwise()
	assert.Equal(t, true, ok, "Path.IsClockwise failed to detect CW path")
}

func TestFlatten(t *testing.T) {
	p1 := path(a, b)
	p2 := path(b, c, d)
	path := Path{}
	if ok := path.Append(&p1); !ok {
		t.Fatal("should be ok")
	}
	if ok := path.Append(&p2); !ok {
		t.Fatal("should be ok")
	}

	path.Flatten()
	pts := path.Points()
	expPts := []v.Vector{a, b, c, d}

	assert.Equal(t, expPts, pts, "Path.Flatten failed")
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
