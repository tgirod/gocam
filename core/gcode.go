package core

import (
	"fmt"

	"github.com/joushou/gocnc/gcode"
)

// this file contains the code to export gcode from the internal representation

func (p *Path) Gcode() []gcode.Block {
	var blocks = make([]gcode.Block, 0, p.Len()+2)

	// add comment
	comm := &gcode.Comment{fmt.Sprintf("Block-name: %s", p.Name), false}
	blocks = append(blocks, gcode.Block{[]gcode.Node{comm}, false})

	// convert moves to gcode

	// initial G0 move to the starting point
	g0 := gcode.Block{}
	g0.AppendNode(&gcode.Word{'G', 0})
	g0.AppendNode(&gcode.Word{'X', p.Start().X})
	g0.AppendNode(&gcode.Word{'Y', p.Start().Y})
	blocks = append(blocks, g0)

	// subsequent G1 moves
	for _, l := range p.Lines {
		start := l.Start()
		end := l.End()

		b := gcode.Block{}

		if l.Bulge == 0 {
			// straight line
			b.AppendNode(&gcode.Word{'G', 1})
			b.AppendNode(&gcode.Word{'X', end.X})
			b.AppendNode(&gcode.Word{'Y', end.Y})
		} else {
			if l.Bulge > 0 {
				// CCW arc
				b.AppendNode(&gcode.Word{'G', 3})
			} else {
				// CW arc
				b.AppendNode(&gcode.Word{'G', 2})
			}
			// arc's endpoint
			b.AppendNode(&gcode.Word{'X', end.X})
			b.AppendNode(&gcode.Word{'Y', end.Y})
			// center (absolute)
			c, _, _, _ := BulgeToArc(start, end, l.Bulge)
			// center (relative to the start)
			c = c.Diff(start)

			b.AppendNode(&gcode.Word{'I', c.X})
			b.AppendNode(&gcode.Word{'J', c.Y})
		}
		blocks = append(blocks, b)
	}
	return blocks
}

func (m *Model) Gcode() *gcode.Document {
	doc := new(gcode.Document)
	for _, path := range m.Paths {
		doc.Blocks = append(doc.Blocks, path.Gcode()...)
	}
	return doc
}
