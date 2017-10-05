package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rpaloschi/dxf-go/document"
	"github.com/rpaloschi/dxf-go/entities"
)

func main() {
	fname := os.Args[1]
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	dxf, err := document.DxfDocumentFromStream(file)
	if err != nil {
		log.Fatal(err)
	}

	doc := Document{}

	for _, entity := range dxf.Entities.Entities {
		if line, ok := entity.(*entities.Line); ok {
			doc.AppendPath(*Line(line))
		}
		if polyline, ok := entity.(*entities.Polyline); ok {
			doc.AppendPath(*Polyline(polyline))
		}
		if lwpolyline, ok := entity.(*entities.LWPolyline); ok {
			doc.AppendPath(*LWPolyline(lwpolyline))
		}
		if arc, ok := entity.(*entities.Arc); ok {
			doc.AppendPath(*Arc(arc))
		}
	}

	doc.Regroup()
	fmt.Println(doc.Export(5))
}
