package main

// This file contains the code to import entities from the DXF format to the
// internal representation of the program

import (
	"io"
	"math"

	v "github.com/joushou/gocnc/vector"
	"github.com/rpaloschi/dxf-go/core"
	"github.com/rpaloschi/dxf-go/document"
	"github.com/rpaloschi/dxf-go/entities"
)

type Importer struct {
	Precision float64
	Ignored   int // number of ignored entities during import
	Discarded int // number of discarded entities (duplicates)
}

func NewImporter() *Importer {
	return &Importer{
		Precision: 1E2,
		Ignored:   0,
		Discarded: 0,
	}
}

func (im *Importer) Import(stream io.Reader) (*Model, error) {
	doc, err := document.DxfDocumentFromStream(stream)
	if err != nil {
		return nil, err
	}

	mod := &Model{}

	Log.Println("Importing entities")
	for _, e := range doc.Entities.Entities {
		move := im.ImportEntity(e)
		mod.Append(move)
	}

	Log.Println("Imported entities: ", len(*mod))
	Log.Println("Ignored entities:  ", im.Ignored)
	Log.Println("Discarded entities:", im.Discarded)

	return mod, nil
}

func (im *Importer) ImportPoint(p core.Point) v.Vector {
	x := math.Floor(p.X*im.Precision) / im.Precision
	y := math.Floor(p.Y*im.Precision) / im.Precision
	z := math.Floor(p.Z*im.Precision) / im.Precision
	return v.Vector{x, y, z}
}

func (im *Importer) ImportEntity(e entities.Entity) Move {
	switch e := e.(type) {
	case *entities.Line:
		return im.ImportLine(e)
	case *entities.Polyline:
		return im.ImportPolyline(e)
	case *entities.LWPolyline:
		return im.ImportLWPolyline(e)
	case *entities.Arc:
		return im.ImportArc(e)
	case *entities.Circle:
		return im.ImportCircle(e)
	default:
		Log.Printf("Ignored entity %T\n", e)
		im.Ignored++
		return nil
	}
}

func (im *Importer) ImportLine(e *entities.Line) *Line {
	from := im.ImportPoint(e.Start)
	to := im.ImportPoint(e.End)
	return &Line{from, to}
}

func (im *Importer) ImportPolyline(e *entities.Polyline) Path {
	p := make(Path, 0, len(e.Vertices)-1)
	for i := 0; i < len(e.Vertices)-1; i++ {
		from := im.ImportPoint(e.Vertices[i].Location)
		to := im.ImportPoint(e.Vertices[i+1].Location)
		p = append(p, &Line{from, to})
	}
	return p
}

func (im *Importer) ImportLWPolyline(e *entities.LWPolyline) Path {
	pts := e.Points
	p := make(Path, 0, len(pts)-1)
	for i := 0; i < len(pts)-1; i++ {
		from := im.ImportPoint(pts[i].Point)
		to := im.ImportPoint(pts[i+1].Point)
		bulge := pts[i].Bulge

		if bulge == 0 {
			p = append(p, &Line{from, to})
		} else {
			center, radius, startAngle, endAngle := bulgeToArc(from, to, bulge)
			if endAngle < startAngle {
				endAngle += math.Pi * 2
			}
			startPoint := pol2car(startAngle, radius).Sum(center)
			endPoint := pol2car(endAngle, radius).Sum(center)
			p = append(p, &Arc{startPoint, endPoint, center, false})
		}
	}
	return p
}

func (im *Importer) ImportArc(e *entities.Arc) *Arc {
	center := im.ImportPoint(e.Center)
	startAngle := deg2rad(e.StartAngle)
	endAngle := deg2rad(e.EndAngle)
	if endAngle < startAngle {
		endAngle += math.Pi * 2
	}
	radius := e.Radius
	startPoint := pol2car(startAngle, radius).Sum(center)
	endPoint := pol2car(endAngle, radius).Sum(center)
	return &Arc{startPoint, endPoint, center, false}
}

// import a circle as two 180 degrees arcs
func (im *Importer) ImportCircle(e *entities.Circle) Path {
	center := im.ImportPoint(e.Center)
	radius := e.Radius
	a := center.Sum(v.Vector{radius, 0, 0})
	b := center.Sum(v.Vector{-radius, 0, 0})
	return Path{
		&Arc{a, b, center, false},
		&Arc{b, a, center, false},
	}
}
