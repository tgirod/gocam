package core

import (
	"fmt"
	"strings"
)

// Model is a list of paths
type Model struct {
	Paths []Path
}

// Len returns the number of paths in the model
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

// Append adds a path to the model
func (mod *Model) Append(p *Path) {
	mod.Paths = append(mod.Paths, *p)
}

// Remove deletes the path at index idx from the model and returns it
func (mod *Model) Remove(idx int) *Path {
	p := &mod.Paths[idx]
	mod.Paths = append(mod.Paths[:idx], mod.Paths[idx+1:]...)
	return p
}

// JoinPath adds a path to the model. Instead of appending the path to the
// model, this method is searching for open prepending and appending paths in
// order to join the new path to existing ones
func (mod *Model) JoinPath(path *Path) {
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
			if !cur.IsClosed() && cur.EndPoint().Equals(path.StartPoint()) {
				// found prepending path
				Log.Printf("\tfound prepending path: %s\n", cur.Name)
				pre = i
			}
			if !cur.IsClosed() && path.EndPoint().Equals(cur.StartPoint()) {
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
