package main

import (
	"fmt"

	"github.com/joushou/gocnc/gcode"
)

// this file contains the code to export gcode from the internal representation

func (p *Path) Gcode() []gcode.Block {
	var blocks = make([]gcode.Block, p.Len()+2)

	// add comment
	comm := &gcode.Comment{fmt.Sprintf("Block-name: %s", p.Handle), false}
	blocks[0] = gcode.Block{[]gcode.Node{comm}, false}
	// initial G0 move
	blocks[1] = gcode.Block{}
	blocks[1].AppendNode(&gcode.Word{'G', 0})
	blocks[1].AppendNode(&gcode.Word{'X', p.Start.X})
	blocks[1].AppendNode(&gcode.Word{'Y', p.Start.Y})

	var start, end Vec
	// convert moves to gcode
	for i := 0; i < p.Len(); i++ {
		if i == 0 {
			start = p.Start
		} else {
			start = p.Moves[i-1].End()
		}
		end = p.Moves[i].End()

		b := gcode.Block{}
		switch move := p.Moves[i].(type) {
		default:
			panic("this should not happen")
		case *Line:
			if move.Bulge == 0 {
				// straight line
				b.AppendNode(&gcode.Word{'G', 1})
				b.AppendNode(&gcode.Word{'X', end.X})
				b.AppendNode(&gcode.Word{'Y', end.Y})
			} else {
				if move.Bulge > 0 {
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
				c, _, _, _ := BulgeToArc(start, end, move.Bulge)
				// center (relative to the start)
				c = c.Sub(start)

				b.AppendNode(&gcode.Word{'I', c.X})
				b.AppendNode(&gcode.Word{'J', c.Y})
			}
		}
		blocks[i+2] = b
	}
	return blocks
}

func (d *Document) Gcode() *gcode.Document {
	gd := new(gcode.Document)
	for _, path := range d.Paths {
		gd.Blocks = append(gd.Blocks, path.Gcode()...)
	}
	return gd
}
