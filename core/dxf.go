package core

// This file contains the code to import entities from the DXF format to the
// internal representation of the program

import (
	"io"
	"math"

	"github.com/rpaloschi/dxf-go/core"
	"github.com/rpaloschi/dxf-go/document"
	"github.com/rpaloschi/dxf-go/entities"
)

// ImportDXF a stream to a DXF file and returns a model or an error if the
// import process went wrong
func ImportDXF(stream io.Reader) (*Model, error) {
	dxf, err := document.DxfDocumentFromStream(stream)
	if err != nil {
		return nil, err
	}

	mod := new(Model)
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

		mod.JoinPath(path)
	}
	return mod, nil
}

// NewVector converts a Point object to a Vector
func NewVector(p core.Point) Vector {
	return Vector{p.X, p.Y}
}

// NewLine converts a Line entity to a path
func NewLine(e *entities.Line) *Path {
	start := NewVector(e.Start)
	end := NewVector(e.End)
	path := NewPath(e.Handle)
	path.Append(Line{start, end, 0})
	return path
}

// NewPolyline converts a polyline entity to a path
func NewPolyline(e *entities.Polyline) *Path {
	path := NewPath(e.Handle)
	// move along vertices
	start := NewVector(e.Vertices[0].Location)
	for _, v := range e.Vertices[1:] {
		end := NewVector(v.Location)
		path.Append(Line{start, end, 0})
		start = end
	}
	return path
}

// NewLWPolyline converts a LWpolyline entity to a path
func NewLWPolyline(e *entities.LWPolyline) *Path {
	path := NewPath(e.Handle)
	// move along vertices
	start := NewVector(e.Points[0].Point)
	for _, p := range e.Points[1:] {
		end := NewVector(p.Point)
		bulge := p.Bulge
		path.Append(Line{start, end, bulge})
		start = end
	}
	return path
}

// NewArc converts a Arc entity to a path. The arc is converted to line+bulge format
func NewArc(e *entities.Arc) *Path {
	center := NewVector(e.Center)
	startAngle := Radians(e.StartAngle)
	endAngle := Radians(e.EndAngle)
	if endAngle < startAngle {
		endAngle += math.Pi * 2
	}
	radius := e.Radius

	startPoint, endPoint, bulge := ArcToBulge(center, radius, startAngle, endAngle)

	path := NewPath(e.Handle)
	path.Append(Line{startPoint, endPoint, bulge})
	return path
}
