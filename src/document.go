package main

import (
	"fmt"
	"strings"
)

type Document struct {
	Paths []Path
}

func (doc *Document) AppendPath(p Path) {
	doc.Paths = append(doc.Paths, p)
}

func (doc *Document) RemovePath(idx int) {
	if idx < len(doc.Paths) {
		doc.Paths = append(doc.Paths[:idx], doc.Paths[idx+1:]...)
	}
}

func (doc *Document) Regroup() {
	for i := 0; i < len(doc.Paths); i++ {
		//fmt.Printf("Number of paths: %d\n", len(doc.Paths))
		cur := &doc.Paths[i]
		if !cur.IsClosed() {
			end := cur.End()
			//fmt.Printf("searching for a follower for %s starting at %v\n", cur.Handle, end)
			// search for a path that starts at "end" and join it to "cur"
			for j, p := range doc.Paths {
				if !p.IsClosed() && p.Start().Equals(end) {
					//fmt.Printf("found %s, starting at %v\n", p.Handle, p.Start())
					// joining paths
					cur.AppendBlocks(p.Blocks...)
					doc.RemovePath(j)
					// a path has been joined to cur. Let's try again
					i--
					break
				}
			}
		}
	}
	for i, _ := range doc.Paths {
		doc.Paths[i].Handle = fmt.Sprintf("%d", i)
	}
}

func (doc *Document) Export(precision int) string {
	l := make([]string, len(doc.Paths))
	for idx, p := range doc.Paths {
		l[idx] = p.Export(precision)
	}
	return strings.Join(l, "\n")
}
