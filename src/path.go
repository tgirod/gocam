package main

import (
	"fmt"
	"strings"
)

type Path struct {
	Blocks []Block // a slice of blocks
	start  Point   // starting point of the path. Starts with a G0 move there
	Handle string  // a name taken from the DXF file
}

// Creates a new Path
func NewPath(start Point, handle string) *Path {
	path := new(Path)
	path.start = start
	path.Handle = handle
	return path
}

// returns the starting coordinates of the path
func (p *Path) Start() Point {
	return p.start
}

// return the coordinates of the path's end
func (p *Path) End() Point {
	if len(p.Blocks) == 0 {
		panic("trying to access last element of an empty slice")
	}
	return p.Blocks[len(p.Blocks)-1].End()
}

func (p *Path) IsClosed() bool {
	return p.start.Equals(p.End())
}

// add a block to the end of the path
func (p *Path) AppendBlock(b Block) {
	p.Blocks = append(p.Blocks, b)
}

// add multiple blocks to the path
func (p *Path) AppendBlocks(bs ...Block) {
	for _, b := range bs {
		p.AppendBlock(b)
	}
}

// exports gcode
func (p *Path) Export(precision int) string {
	l := make([]string, len(p.Blocks)+2)
	l[0] = fmt.Sprintf("(Block-name: %s)", p.Handle)
	g0 := *newMove(0, p.start)
	l[1] = g0.Export(precision)
	for idx, b := range p.Blocks {
		l[idx+2] = b.Export(precision)
	}
	return strings.Join(l, "\n")
}

// reverse path
func (p *Path) Reverse() {
	// reverse the slice
	for i, j := 0, len(p.Blocks)-1; i < j; i, j = i+1, j-1 {
		p.Blocks[i], p.Blocks[j] = p.Blocks[j], p.Blocks[i]
	}
}
