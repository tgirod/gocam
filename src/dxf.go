package main

// This file contains the code to import entities from the DXF format to the
// internal representation of the program

import (
	"math"

	"github.com/rpaloschi/dxf-go/core"
	"github.com/rpaloschi/dxf-go/document"
	"github.com/rpaloschi/dxf-go/entities"
)

func NewVec(p core.Point) Vec {
	return Vec{p.X, p.Y}
}

func NewLine(e *entities.Line) *Path {
	start := NewVec(e.Start)
	end := NewVec(e.End)
	path := NewPath(end, e.Handle)
	path.AppendMove(&Line{start, 0})
	return path
}

func NewPolyline(e *entities.Polyline) *Path {
	start := NewVec(e.Vertices[0].Location)
	path := NewPath(start, e.Handle)
	// move along vertices
	for _, v := range e.Vertices[1:] {
		end := NewVec(v.Location)
		path.AppendMove(&Line{end, 0})
	}
	return path
}

func NewLWPolyline(e *entities.LWPolyline) *Path {
	start := NewVec(e.Points[0].Point)
	path := NewPath(start, e.Handle)
	// move along vertices
	for _, p := range e.Points[1:] {
		end := NewVec(p.Point)
		bulge := p.Bulge
		path.AppendMove(&Line{end, bulge})
	}
	return path
}

func NewArc(e *entities.Arc) *Path {
	center := NewVec(e.Center)
	startAngle := Radians(e.StartAngle)
	endAngle := Radians(e.EndAngle)
	if endAngle < startAngle {
		endAngle += math.Pi * 2
	}
	radius := e.Radius

	startPoint, endPoint, bulge := ArcToBulge(center, radius, startAngle, endAngle)

	path := NewPath(startPoint, e.Handle)
	path.AppendMove(&Line{endPoint, bulge})
	return path
}

func NewDocumentFromDxf(dxf *document.DxfDocument) *Document {
	doc := new(Document)
	// import each entity as a separate path
	for _, entity := range dxf.Entities.Entities {
		var path *Path
		switch e := entity.(type) {
		default:
			//fmt.Printf("unsupported type %T\n", e)
		case *entities.Line:
			//fmt.Printf("importing %T %s\n", e, e.Handle)
			path = NewLine(e)
		case *entities.Polyline:
			//fmt.Printf("importing %T %s\n", e, e.Handle)
			path = NewPolyline(e)
		case *entities.LWPolyline:
			//fmt.Printf("importing %T %s\n", e, e.Handle)
			path = NewLWPolyline(e)
		case *entities.Arc:
			//fmt.Printf("importing %T %s\n", e, e.Handle)
			path = NewArc(e)
		}

		doc.JoinPath(path)
	}
	return doc
}
