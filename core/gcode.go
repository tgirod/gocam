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

func move(v Vector) gcode.Block {
	m := gcode.Block{}
	m.AppendNodes(
		word('G', 0),
		word('X', v.X),
		word('Y', v.Y))
	return m
}

// Gcode generates gcode for the given path
func (p *Path) Gcode() []gcode.Block {
	var blocks = make([]gcode.Block, 0, p.Len()+2)

	// add a header (bCNC compatible, for easy visualisation)
	blocks = append(blocks, header(p.Name))

	// initial G0 move to the starting point
	blocks = append(blocks, move(p.StartPoint()))
	prev := p.StartPoint() // where last move ended

	// subsequent G1 moves
	for _, m := range p.Moves {
		// add G0 move if the line does not start where the previous one ended
		if !m.StartPoint().Equals(prev) {
			Log.Printf("Path is not continuous. Adding a G0 move from %v to %v", prev, m.StartPoint())
			blocks = append(blocks, move(m.StartPoint()))
		}
		// add actual move (G1, G2 or G3)
		b := gcode.Block{}
		switch t := m.(type) {
		default:
			Log.Printf("Unexpected type %T", t)
		case *Line:
			b.AppendNodes(
				word('G', 1),
				word('X', t.End.X),
				word('Y', t.End.Y))
		case *Arc:
			if t.CW {
				b.AppendNode(word('G', 2))
			} else {
				b.AppendNode(word('G', 3))
			}
			// arc's endpoint
			b.AppendNode(word('X', t.End.X))
			b.AppendNode(word('Y', t.End.Y))
			// center (relative to the start)
			center := t.Center.Sub(t.Start)
			b.AppendNode(word('I', center.X))
			b.AppendNode(word('J', center.Y))
		}
		blocks = append(blocks, b)
		prev = m.EndPoint()
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
