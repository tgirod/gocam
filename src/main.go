package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rpaloschi/dxf-go/document"
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

	doc := NewDocumentFromDxf(dxf)
	//doc.Regroup()
	fmt.Println(doc.Gcode().Export(5))
}
