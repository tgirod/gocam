package main

import "github.com/joushou/gocnc/gcode"

type Block struct {
	gcode.Block        // the actual gcode stuff
	end         *Point // copy of XYZ words if available
}

// set XYZ words for the block
func (b *Block) SetEnd(p Point) {
	b.RemoveAddress('X', 'Y', 'Z')
	b.AppendNodes(
		&gcode.Word{'X', p.X},
		&gcode.Word{'Y', p.Y},
		&gcode.Word{'Z', p.Z})
	b.end = &p
}

func (b *Block) End() Point {
	return *b.end
}

// create a move block base with G + XYZ words
func newMove(code float64, p Point) *Block {
	g := new(Block)
	g.AppendNode(&gcode.Word{'G', code})
	g.SetEnd(p)
	return g
}
