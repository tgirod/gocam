package main

import (
	"github.com/joushou/gocnc/gcode"
	"github.com/rpaloschi/dxf-go/entities"
)

func Line(l *entities.Line) *Path {
	path := NewPath(Point{l.Start}, l.Handle)
	g1 := newMove(1, Point{l.End})
	path.AppendBlock(*g1)
	return path
}

func Polyline(p *entities.Polyline) *Path {
	start := Point{p.Vertices[0].Location}
	path := NewPath(start, p.Handle)
	// move along vertices
	for _, v := range p.Vertices[1:] {
		g1 := newMove(1, Point{v.Location})
		path.AppendBlock(*g1)
	}
	return path
}

func LWPolyline(lwp *entities.LWPolyline) *Path {
	start := Point{lwp.Points[0].Point}
	path := NewPath(start, lwp.Handle)
	// move along vertices
	for _, p := range lwp.Points[1:] {
		g1 := newMove(1, Point{p.Point})
		path.AppendBlock(*g1)
	}
	return path
}

func Arc(a *entities.Arc) *Path {
	center := Point{a.Center}
	start := pol2car(center, a.Radius, a.StartAngle)
	end := pol2car(center, a.Radius, a.EndAngle)

	path := NewPath(start, a.Handle)
	// CCW arc
	g3 := newMove(3, end)
	// add center (relative to the starting point
	i := center.X - start.X
	j := center.Y - start.Y
	g3.AppendNodes(
		&gcode.Word{'I', i},
		&gcode.Word{'J', j})
	path.AppendBlock(*g3)
	return path
}
