package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tgirod/gocam/dxf"
	"github.com/tgirod/gocam/gcode"
)

func main() {
	fname := os.Args[1]
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	if mod, err := dxf.Import(file); err == nil {
		gc := gcode.Export(mod)
		fmt.Println(gc.Export(5))
	} else {
		log.Fatal(err)
	}
}
