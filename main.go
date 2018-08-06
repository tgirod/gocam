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
	var mod *Model
	if mod, err = im.Import(file); err != nil {
		Log.Fatal(err)
	}

	doc := mod.Gcode()
	fmt.Println(doc.Export(2))
}
