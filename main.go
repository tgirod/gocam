package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tgirod/gocam/core"
)

func main() {
	fname := os.Args[1]
	file, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	if mod, err := core.ImportDXF(file); err == nil {
		for i := 0; i < mod.Len(); i++ {
			mod.Paths[i].StartFrom(core.Vector{0, 0})
		}
		fmt.Println(mod.Gcode().Export(5))
	} else {
		log.Fatal(err)
	}
}
