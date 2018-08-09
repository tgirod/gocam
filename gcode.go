package main

import (
	"fmt"

	"github.com/joushou/gocnc/gcode"
)

type Gcoder interface {
	Gcode() gcode.Block
}

func word(address rune, command float64) gcode.Node {
	return &gcode.Word{
		Address: address,
		Command: command}
}

func header(id int) gcode.Node {
	return &gcode.Comment{
		Content: fmt.Sprintf("Block-name: %d", id),
		EOL:     false,
	}
}

func xy(v Vector) []gcode.Node {
	return []gcode.Node{
		word('X', v.X),
		word('Y', v.Y),
	}
}

func ij(v Vector) []gcode.Node {
	return []gcode.Node{
		word('I', v.X),
		word('J', v.Y),
	}
}

func move(v Vector) gcode.Block {
	return gcode.Block{
		Nodes: []gcode.Node{
			word('G', 0),
			word('X', v.X),
			word('Y', v.Y),
		},
	}
}
