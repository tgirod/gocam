package core

import (
	"fmt"
	"strings"
)

// A model is a list of paths
type Model struct {
	Paths []Path
}

func (mod *Model) Len() int {
	return len(mod.Paths)
}

func (mod *Model) String() string {
	l := make([]string, mod.Len())
	for i, m := range mod.Paths {
		l[i] = fmt.Sprintf("%s", &m)
	}
	return strings.Join(l, "\n")
}

func (mod *Model) Append(p *Path) {
	mod.Paths = append(mod.Paths, *p)
}

func (mod *Model) Remove(idx int) *Path {
	p := &mod.Paths[idx]
	mod.Paths = append(mod.Paths[:idx], mod.Paths[idx+1:]...)
	return p
}

// Add a path to the model. This method will try to join the new path to the
// existing ones by looking for possible prepending and appending paths. Paths
// will be joined as necessary.
func (mod *Model) AddPath(path *Path) {
	Log.Printf("adding path %s\n", path.Name)
	if path.IsClosed() {
		if path.IsClockwise() {
			path.Reverse()
		}
		mod.Append(path)
	} else {
		// search for prepending and appending paths
		pre := -1
		post := -1

		for i := 0; i < mod.Len() && (pre == -1 || post == -1); i++ {
			cur := &mod.Paths[i]
			if !cur.IsClosed() && cur.End().Equals(path.Start()) {
				// found prepending path
				Log.Printf("\tfound prepending path: %s\n", cur.Name)
				pre = i
			}
			if !cur.IsClosed() && path.End().Equals(cur.Start()) {
				// found appending path
				Log.Printf("\tfound appending path: %s\n", cur.Name)
				post = i
			}
		}

		preHandle := "xxx"
		postHandle := "xxx"
		if pre == -1 && post == -1 {
			// isolated path
			// xxx -> PATH -> xxx
			mod.Append(path)
		} else if post == -1 {
			// prepending path only
			// PRE -> PATH
			preHandle = mod.Paths[pre].Name
			mod.Paths[pre].Join(path)
		} else if pre == -1 {
			// appending path only
			// PATH -> POST
			postHandle = mod.Paths[post].Name
			path.Join(&mod.Paths[post])
			mod.Paths[post] = *path
		} else if pre == post {
			// loop-closing path
			// PREPOST -> PATH
			preHandle = mod.Paths[pre].Name
			postHandle = mod.Paths[post].Name
			mod.Paths[pre].Join(path)
			if mod.Paths[pre].IsClockwise() {
				mod.Paths[pre].Reverse()
			}
		} else {
			// prepending and appending paths
			// PRE -> PATH -> POST
			preHandle = mod.Paths[pre].Name
			postHandle = mod.Paths[post].Name
			mod.Paths[pre].Join(path)
			mod.Paths[pre].Join(&mod.Paths[post])
			mod.Remove(post)
		}
		Log.Printf("\t%s->%s->%s\n", preHandle, path.Name, postHandle)
	}
}
