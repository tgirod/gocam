package model

import (
	"testing"
)

var a, b, c, d = Vector{-1, -1}, Vector{1, -1}, Vector{1, 1}, Vector{-1, 1}

func poly(pts ...Vector) *Path {
	p := Path{}
	for i := 0; i < len(pts)-1; i++ {
		p.Moves = append(p.Moves, &Line{pts[i], pts[i+1]})
	}
	return &p
}

func TestEqual(t *testing.T) {
	p1 := poly(a, b, c)
	p2 := poly(a, b, c)
	if !p1.Equals(p2) {
		t.Error("PathEquals failed")
	}
}

func TestStart(t *testing.T) {
	p := poly(a, b, c)
	if p.Start() != a {
		t.Error("Start failed")
	}
}

func TestEnd(t *testing.T) {
	p := poly(a, b, c)
	if p.End() != c {
		t.Error("End failed")
	}
}

func TestIsContinuous(t *testing.T) {
	p := poly(a, b, c, d)
	if !p.IsContinuous() {
		t.Error("Path.IsContinuous failed: path is continuous")
	}
	p.Append(&Line{a, b})
	if p.IsContinuous() {
		t.Error("Path.IsContinuous failed: path is not continuous")
	}
}

func TestReverse(t *testing.T) {
	p1 := poly(a, b, c, a)
	p2 := poly(a, c, b, a)
	p2.Reverse()

	if !p1.Equals(p2) {
		t.Error("Path.Reverse failed")
	}
}

func TestPoints(t *testing.T) {
	p := poly(a, b, c)
	pts := p.Points()
	if len(pts) != 3 || pts[0] != a || pts[1] != b || pts[2] != c {
		t.Error("Path.Points failed")
	}
}

func TestAppend(t *testing.T) {
	p := poly(a, b)
	p.Append(&Line{b, c})
	if !p.Equals(poly(a, b, c)) {
		t.Error("Path.Append failed")
	}
}

func TestIsClosed(t *testing.T) {
	p := poly(a, b, c, a)
	if !p.IsClosed() {
		t.Error("Path.IsClosed failed")
	}
}

func TestIsClockwise(t *testing.T) {
	p := poly(a, b, c, a)
	if p.IsClockwise() {
		t.Error("Path.IsClockwise failed to detect CCW path")
	}
}

func TestIsClockwise2(t *testing.T) {
	p := poly(a, c, b, a)
	if !p.IsClockwise() {
		t.Error("Path.IsClockwise failed to detect CW path")
	}
}

func TestHasInnerLoop(t *testing.T) {
	var data = []struct {
		in        *Path
		hasLoop   bool
		startLoop int
		endLoop   int
		doc       string
	}{
		{poly(a, b, c, d), false, 0, 0, "no loop"},
		{poly(a, b, c, d, a), false, 0, 0, "outer loop"},
		{poly(a, b, c, a, d), true, 0, 3, "start loop"},
		{poly(a, b, c, d, b), true, 1, 4, "end loop"},
		{poly(a, b, c, b, d), true, 1, 3, "inner loop"},
	}
	for _, d := range data {
		hasLoop, startLoop, endLoop := d.in.HasInnerLoop()
		if d.hasLoop != hasLoop || d.startLoop != startLoop || d.endLoop != endLoop {
			t.Error("Path.HasInnerLoop failed", d.doc, hasLoop, startLoop, endLoop)
		}
	}
}

func TestSplitInnerLoops(t *testing.T) {
	var data = []struct {
		in  *Path
		out []*Path
		doc string
	}{
		{poly(a, b, c, d), []*Path{poly(a, b, c, d)}, "no loop"},
		{poly(a, b, c, d, a), []*Path{poly(a, b, c, d, a)}, "outer loop"},
		{poly(a, b, c, a, d), []*Path{poly(a, b, c, a), poly(a, d)}, "start loop"},
		{poly(a, b, c, d, b), []*Path{poly(b, c, d, b), poly(a, b)}, "end loop"},
		{poly(a, b, c, b, d), []*Path{poly(b, c, b), poly(a, b, d)}, "inner loop"},
	}
	for _, d := range data {
		paths := d.in.SplitInnerLoops()
		if len(paths) != len(d.out) {
			t.Error("Path.SplitInnerLoops failed: wrong number of loops")
		}
		for i := 0; i < len(paths); i++ {
			if !paths[i].Equals(d.out[i]) {
				t.Error("Path.SplitInnerLoops failed on", d.doc, "test", i)
				t.Log("RETURNED", paths[i])
				t.Log("EXPECTED", d.out[i])
			}
		}
	}
}
