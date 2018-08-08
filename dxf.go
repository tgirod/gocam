package main

// This file contains the code to import entities from the DXF format to the
// internal representation of the program

import (
	"io"
	"math"

	"github.com/davecgh/go-spew/spew"
	"github.com/rpaloschi/dxf-go/core"
	"github.com/rpaloschi/dxf-go/document"
	"github.com/rpaloschi/dxf-go/entities"
)

type Importer struct {
	Precision int
	Imported  int // number of imported entities
	Ignored   int // number of ignored entities
	Discarded int // number of discarded entities (duplicates)
	Model     *Model
}

func NewImporter() *Importer {
	return &Importer{
		Precision: 3,
		Model:     &Model{},
	}
}

func (im *Importer) Import(stream io.Reader) error {
	doc, err := document.DxfDocumentFromStream(stream)
	if err != nil {
		return err
	}

	Log.Println("Importing entities")
	for _, e := range doc.Entities.Entities {
		im.ImportEntity(e)
	}
	im.Model.Merge()

	Log.Println("Imported entities: ", im.Imported)
	Log.Println("Ignored entities:  ", im.Ignored)
	Log.Println("Discarded entities:", im.Discarded)

	return nil
}

func (im *Importer) ImportPoint(p core.Point) Vector {
	pre := math.Pow10(im.Precision)
	x := math.Floor(p.X*pre) / pre
	y := math.Floor(p.Y*pre) / pre
	z := math.Floor(p.Z*pre) / pre
	return Vector{x, y, z}
}

func (im *Importer) ImportEntity(e entities.Entity) {
	switch e := e.(type) {
	case *entities.Line:
		im.ImportLine(e)
	case *entities.Polyline:
		im.ImportPolyline(e)
	case *entities.LWPolyline:
		im.ImportLWPolyline(e)
	case *entities.Arc:
		im.ImportArc(e)
	case *entities.Circle:
		im.ImportCircle(e)
	case *entities.Spline:
		im.ImportSpline(e)
	default:
		Log.Printf("Ignored entity %T\n", e)
		im.Ignored++
	}
}

func (im *Importer) ImportLine(e *entities.Line) {
	from := im.ImportPoint(e.Start)
	to := im.ImportPoint(e.End)
	im.Model.Append(&Line{from, to})
	im.Imported++
}

func (im *Importer) ImportPolyline(e *entities.Polyline) {
	p := make(Path, 0, len(e.Vertices)-1)
	for i := 0; i < len(e.Vertices)-1; i++ {
		from := im.ImportPoint(e.Vertices[i].Location)
		to := im.ImportPoint(e.Vertices[i+1].Location)
		p = append(p, &Line{from, to})
	}
	im.Model.Append(p)
	im.Imported++
}

func (im *Importer) ImportLWPolyline(e *entities.LWPolyline) {
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
	im.Model.Append(p)
	im.Imported++
}

func (im *Importer) ImportArc(e *entities.Arc) {
	center := im.ImportPoint(e.Center)
	startAngle := deg2rad(e.StartAngle)
	endAngle := deg2rad(e.EndAngle)
	if endAngle < startAngle {
		endAngle += math.Pi * 2
	}
	radius := e.Radius
	startPoint := pol2car(startAngle, radius).Sum(center)
	endPoint := pol2car(endAngle, radius).Sum(center)
	im.Model.Append(&Arc{startPoint, endPoint, center, false})
	im.Imported++
}

// import a circle as two 180 degrees arcs
func (im *Importer) ImportCircle(e *entities.Circle) {
	center := im.ImportPoint(e.Center)
	radius := e.Radius
	a := center.Sum(Vector{radius, 0, 0})
	b := center.Sum(Vector{-radius, 0, 0})
	p := Path{
		&Arc{a, b, center, false},
		&Arc{b, a, center, false},
	}
	im.Model.Append(p)
	im.Imported++
}

func (im *Importer) ImportSpline(e *entities.Spline) {
	// FIXME
	s := &Spline{}
	s.Degree = e.Degree
	s.Closed = e.Closed
	s.Knots = e.KnotValues
	s.Controls = make([]Vector, len(e.ControlPoints))
	for i, p := range e.ControlPoints {
		s.Controls[i] = im.ImportPoint(p)
	}
	if len(e.Weights) != 0 {
		s.Weights = e.Weights
	} else {
		s.Weights = make([]float64, len(s.Controls))
		for i := range s.Controls {
			s.Weights[i] = 1
		}
	}
	// FIXME
	// im.Model.Append(s)
	// im.Imported++
	spew.Dump(s)
	// max := float64(s.Knots[len(s.Knots)-1])
	// for i := 0; i < 101; i++ {
	// 	u := float64(i) / 100 * max
	// 	s.eval(u)
	// 	fmt.Println("u", u)
	// 	// fmt.Printf("%f %f\n", v.X, v.Y)
	// }
}
