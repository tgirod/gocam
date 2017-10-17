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
	doc.Regroup()
	//fmt.Println(doc)
	for _, p := range doc.Paths {
		fmt.Printf("%s: closed:%t, clockwise:%t\n", p.Handle, p.IsClosed(), p.IsClockwise())
	}

	//fmt.Println(doc.Gcode().Export(5))
}

//func main() {
//fmt.Println(ArcToBulge(Vec{0, 0}, 100, 0, Radians(180)))
//fmt.Println(BulgeToArc(Vec{100, 0}, Vec{-100, 0}, 1))
//}
