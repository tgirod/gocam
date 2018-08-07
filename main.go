package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fname := os.Args[1]
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	im := NewImporter()
	if err = im.Import(file); err != nil {
		Log.Fatal(err)
	}

	doc := im.Model.Gcode()
	fmt.Println(doc.Export(im.Precision))
}
