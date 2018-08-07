package main

import (
	"github.com/joushou/gocnc/gcode"
)

type Model []Path

// Append Move to the model by concatenating it to an existing one if possible.
func (m *Model) Append(mo Move) {
	for i := 0; i < len(*m); i++ {
		p := &(*m)[i]
		if ok := p.Append(mo); ok {
			return
		}
	}
	p := &Path{}
	p.Append(mo)
	*m = append(*m, *p)
}

func (m *Model) Merge() {
	l := len(*m)
	for i := 0; i < l; i++ {
		// extract the first path
		p := (*m)[0]
		*m = (*m)[1:]
		// re-append it
		m.Append(p)
	}
}

func (m *Model) Gcode() gcode.Document {
	doc := &gcode.Document{}
	for _, p := range *m {
		bs := p.Gcode()
		doc.Blocks = append(doc.Blocks, bs...)
	}
	return *doc
}
