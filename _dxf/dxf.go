package dxf

// This file contains the code to import entities from the DXF format to the
// internal representation of the program

import (
	"io"
	"math"

	"github.com/rpaloschi/dxf-go/core"
	"github.com/rpaloschi/dxf-go/document"
	"github.com/rpaloschi/dxf-go/entities"
	"github.com/tgirod/gocam/model"
	"github.com/tgirod/gocam/util"
)

const PRECISION float64 = 1E-3

func Import(stream io.Reader) (*model.Model, error) {
	doc, err := document.DxfDocumentFromStream(stream)
	if err != nil {
		return nil, err
	}

	util.Log.Println("Importing entities")
	im := NewImporter()

	for _, e := range doc.Entities.Entities {
		im.ImportEntity(e)
	}

	util.Log.Println("Imported entities: ", len(im.Moves))
	util.Log.Println("Ignored entities:  ", im.Ignored)
	util.Log.Println("Discarded entities:", im.Discarded)

	mod := model.New()

	util.Log.Println("Building the model")
	for _, m := range im.Moves {
		mod.Add(model.NewPath(m))
	}

	return mod, nil
}

type Importer struct {
	Points    map[model.Vector]bool // points found in entities
	Moves     []model.Move          // every moves already imported
	Ignored   int                   // number of ignored entities during import
	Discarded int                   // number of discarded entities (duplicates)
}

func NewImporter() *Importer {
	im := new(Importer)
	im.Points = make(map[model.Vector]bool)
	im.Moves = []model.Move{}
	return im
}

func (im *Importer) Add(m model.Move) {
	discard := false
	// discard null move
	if m.Start() == m.End() {
		discard = true
	} else {

		// dicard copy or reverse copy of existing move
		rm := m.Copy()
		rm.Reverse()
		for i := 0; i < len(im.Moves) && !discard; i++ {
			discard = m.Equals(im.Moves[i]) || rm.Equals(im.Moves[i])
		}
	}

	if !discard {
		im.Moves = append(im.Moves, m)
	} else {
		im.Discarded++
	}
}

func (im *Importer) ImportPoint(p core.Point) model.Vector {
	v := model.Vector{p.X, p.Y}
	return im.ImportVector(v)
}

func (im *Importer) ImportVector(v model.Vector) model.Vector {
	if _, ok := im.Points[v]; ok {
		// v already exists in im.Points
		return v
	} else {
		for w, _ := range im.Points {
			d2 := math.Pow(w.X-v.X, 2) + math.Pow(w.Y-v.Y, 2)
			if d2 < PRECISION {
				// w is sufficiently close to v, we ignore v and return w
				return w
			}
		}
		// no point sufficiently close, add v and return
		im.Points[v] = true
		return v
	}
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
	default:
		util.Log.Printf("Ignored entity %T\n", e)
		im.Ignored++
	}
}

func (im *Importer) ImportLine(e *entities.Line) {
	start := im.ImportPoint(e.Start)
	end := im.ImportPoint(e.End)
	im.Add(&model.Line{start, end})
}

func (im *Importer) ImportPolyline(e *entities.Polyline) {
	// move along vertices
	start := im.ImportPoint(e.Vertices[0].Location)
	for _, v := range e.Vertices[1:] {
		end := im.ImportPoint(v.Location)
		im.Add(&model.Line{start, end})
		start = end
	}
}

func (im *Importer) ImportLWPolyline(e *entities.LWPolyline) {
	// move along vertices
	start := im.ImportPoint(e.Points[0].Point)
	for _, p := range e.Points[1:] {
		end := im.ImportPoint(p.Point)
		if p.Bulge == 0 {
			im.Add(&model.Line{start, end})
		} else {
			im.ImportArcFromBulge(start, end, p.Bulge)
		}
		start = end
	}
}

func (im *Importer) ImportArcFromBulge(start model.Vector, end model.Vector, bulge float64) {
	chord := end.Sub(start)                                    // chord of the arc
	theta2 := 2 * math.Atan(bulge)                             // half of included angle
	d := chord.Length() / 2                                    // half of the chord length
	r := d / math.Sin(theta2)                                  // radius
	a := end.Sub(start).Angle()                                // angle of the chord
	center := util.Cartesian(math.Pi/2-theta2+a, r).Add(start) // center
	center = im.ImportVector(center)
	cw := bulge < 0
	im.Add(&model.Arc{start, end, center, cw})
}

func (im *Importer) ImportArc(e *entities.Arc) {
	center := im.ImportPoint(e.Center)
	startAngle := util.Radians(e.StartAngle)
	endAngle := util.Radians(e.EndAngle)
	if endAngle < startAngle {
		endAngle += math.Pi * 2
	}
	radius := e.Radius
	startPoint := im.ImportVector(util.Cartesian(startAngle, radius).Add(center))
	endPoint := im.ImportVector(util.Cartesian(endAngle, radius).Add(center))
	im.Add(&model.Arc{startPoint, endPoint, center, false})
}

// import a circle as two 180 degrees arcs
func (im *Importer) ImportCircle(e *entities.Circle) {
	center := im.ImportPoint(e.Center)
	radius := e.Radius
	a := im.ImportVector(center.Add(model.Vector{radius, 0}))
	b := im.ImportVector(center.Add(model.Vector{-radius, 0}))
	im.Add(&model.Arc{a, b, center, false})
	im.Add(&model.Arc{b, a, center, false})
}
