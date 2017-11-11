package core

import (
	"fmt"

	"github.com/joushou/gocnc/gcode"
)

// this file contains the code to export gcode from the internal representation

func word(address rune, command float64) *gcode.Word {
	return &gcode.Word{
		Address: address,
		Command: command}
}

func header(blockName string) gcode.Block {
	comm := &gcode.Comment{
		Content: fmt.Sprintf("Block-name: %s", blockName),
		EOL:     false}

	block := gcode.Block{}
	block.AppendNode(comm)
	return block
}

func move(x, y float64) gcode.Block {
	m := gcode.Block{}
	m.AppendNodes(
		word('G', 0),
		word('X', x),
		word('Y', y))
	return m
}

// Gcode generates gcode for the given path
func (p *Path) Gcode() []gcode.Block {
	var blocks = make([]gcode.Block, 0, p.Len()+2)

	// add a header (bCNC compatible, for easy visualisation)
	blocks = append(blocks, header(p.Name))

	// initial G0 move to the starting point
	blocks = append(blocks, move(p.Start().X, p.Start().Y))
	prev := p.Start() // where last move ended

	// subsequent G1 moves
	for _, l := range p.Lines {
		// add G0 move if the line does not start where the previous one ended
		if !l.Start.Equals(prev) {
			blocks = append(blocks, move(l.Start.X, l.Start.Y))
		}
		// add actual move (G1, G2 or G3)
		b := gcode.Block{}
		if l.Bulge == 0 {
			// straight line
			b.AppendNodes(
				word('G', 1),
				word('X', l.End.X),
				word('Y', l.End.Y))
		} else {
			if l.Bulge > 0 {
				// CCW arc
				b.AppendNode(word('G', 3))
			} else {
				// CW arc
				b.AppendNode(word('G', 2))
			}
			// arc's endpoint
			b.AppendNode(word('X', l.End.X))
			b.AppendNode(word('Y', l.End.Y))
			// center (absolute)
			c, _, _, _ := BulgeToArc(l.Start, l.End, l.Bulge)
			// center (relative to the start)
			c = c.Sub(l.Start)

			b.AppendNode(word('I', c.X))
			b.AppendNode(word('J', c.Y))
		}
		blocks = append(blocks, b)
		prev = l.End
	}
	return blocks
}

// Gcode generates gcode for the given model m
func (m *Model) Gcode() *gcode.Document {
	doc := new(gcode.Document)
	for _, path := range m.Paths {
		doc.Blocks = append(doc.Blocks, path.Gcode()...)
	}
	return doc
}
