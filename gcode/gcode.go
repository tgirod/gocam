package gcode

import (
	"fmt"

	"github.com/joushou/gocnc/gcode"
	"github.com/tgirod/gocam/model"
	"github.com/tgirod/gocam/util"
)

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

func move(v model.Vector) gcode.Block {
	m := gcode.Block{}
	m.AppendNodes(
		word('G', 0),
		word('X', v.X),
		word('Y', v.Y))
	return m
}

func Export(m *model.Model) gcode.Document {
	doc := gcode.Document{}
	for i, path := range m.Paths {
		name := fmt.Sprintf("%d", i)
		doc.Blocks = append(doc.Blocks, ExportPath(&path, name)...)
	}
	return doc
}

func ExportMove(m model.Move) gcode.Block {
	var b gcode.Block
	switch m := m.(type) {
	case *model.Line:
		b = ExportLine(m)
	case *model.Arc:
		b = ExportArc(m)
	}
	return b
}

func ExportLine(l *model.Line) gcode.Block {
	b := gcode.Block{}
	b.AppendNodes(
		word('G', 1),
		word('X', l.To.X),
		word('Y', l.To.Y))
	return b
}

func ExportArc(a *model.Arc) gcode.Block {
	b := gcode.Block{}
	if a.CW {
		b.AppendNode(word('G', 2))
	} else {
		b.AppendNode(word('G', 3))
	}
	// arc's endpoint
	b.AppendNode(word('X', a.To.X))
	b.AppendNode(word('Y', a.To.Y))
	// center (relative to the start)
	center := a.Center.Sub(a.From)
	b.AppendNode(word('I', center.X))
	b.AppendNode(word('J', center.Y))
	return b
}

// Gcode generates gcode for the given path
func ExportPath(p *model.Path, name string) []gcode.Block {
	var blocks = make([]gcode.Block, 0, p.Len()+2)

	// add a header (bCNC compatible, for easy visualisation)
	blocks = append(blocks, header(name))

	// initial G0 move to the starting point
	blocks = append(blocks, move(p.Start()))
	prev := p.Start() // where last move ended

	// subsequent G1 moves
	for _, m := range p.Moves {
		// add G0 move if the line does not start where the previous one ended
		if m.Start() != prev {
			util.Log.Printf("Path is not continuous. Adding a G0 move from %v to %v", prev, m.Start())
			blocks = append(blocks, move(m.Start()))
		}
		// add actual move (G1, G2 or G3)
		blocks = append(blocks, ExportMove(m))
		prev = m.End()
	}
	return blocks
}
