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
		//core.Log.Println(mod)
		for _, p := range mod.Paths {
			core.Log.Println(p.Name, p.Bounds())
		}
		fmt.Println(mod.Gcode().Export(5))
	} else {
		log.Fatal(err)
	}
}
