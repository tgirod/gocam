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

// Gcode generates gcode for the given path
func (p *Path) Gcode() []gcode.Block {
	var blocks = make([]gcode.Block, 0, p.Len()+2)

	// add a header (bCNC compatible, for easy visualisation)
	blocks = append(blocks, header(p.Name))

	// convert moves to gcode

	// initial G0 move to the starting point
	g0 := gcode.Block{}
	g0.AppendNodes(
		word('G', 0),
		word('X', p.Start().X),
		word('Y', p.Start().Y))
	blocks = append(blocks, g0)

	// subsequent G1 moves
	for _, l := range p.Lines {
		b := gcode.Block{}
		if l.Bulge == 0 {
			// straight line
			b.AppendNodes(
				word('G', 1),
				word('X', l.To.X),
				word('Y', l.To.Y))
		} else {
			if l.Bulge > 0 {
				// CCW arc
				b.AppendNode(word('G', 3))
			} else {
				// CW arc
				b.AppendNode(word('G', 2))
			}
			// arc's endpoint
			b.AppendNode(word('X', l.To.X))
			b.AppendNode(word('Y', l.To.Y))
			// center (absolute)
			c, _, _, _ := BulgeToArc(l.From, l.To, l.Bulge)
			// center (relative to the start)
			c = c.Diff(l.From)

			b.AppendNode(word('I', c.X))
			b.AppendNode(word('J', c.Y))
		}
		blocks = append(blocks, b)
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
