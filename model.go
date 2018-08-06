package main

import (
	"github.com/joushou/gocnc/gcode"
)

type Model []Path

// Append Move to the model by concatenating it to an existing one if possible.
func (m *Model) Append(mo Move) {
	Log.Printf("Appending: %v\n", mo)
	for i := 0; i < len(*m); i++ {
		p := &(*m)[i]
		if ok := p.Append(mo); ok {
			Log.Printf("updated path: %v", p)
			return
		}
	}
	p := &Path{}
	p.Append(mo)
	*m = append(*m, *p)
	Log.Printf("new path: %v", *p)
}

func (m *Model) Gcode() gcode.Document {
	doc := &gcode.Document{}
	for _, p := range *m {
		bs := p.Gcode()
		doc.Blocks = append(doc.Blocks, bs...)
	}
	return *doc
}
