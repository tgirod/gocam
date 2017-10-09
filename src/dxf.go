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
	end := NewVec(e.Vertices[len(e.Vertices)-1].Location)
	path := NewPath(end, e.Handle)
	// move along vertices
	for _, v := range e.Vertices {
		start := NewVec(v.Location)
		path.AppendMove(&Line{start, 0})
	}
	return path
}

func NewLWPolyline(e *entities.LWPolyline) *Path {
	end := NewVec(e.Points[len(e.Points)-1].Point)
	path := NewPath(end, e.Handle)
	// move along vertices
	for _, p := range e.Points {
		start := NewVec(p.Point)
		bulge := p.Bulge
		path.AppendMove(&Line{start, bulge})
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

	path := NewPath(endPoint, e.Handle)
	path.AppendMove(&Line{startPoint, bulge})
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
		doc.Paths = append(doc.Paths, *path)
	}
	return doc
}
