package model

type Path struct {
	Moves []Move
}

func NewPath(m Move) *Path {
	return &Path{[]Move{m}}
}

func (p *Path) Len() int {
	return len(p.Moves)
}

// func (p *Path) String() string {
// 	pts := p.Points()
// 	s := make([]string, len(pts))
// 	for i, pt := range pts {
// 		s[i] = fmt.Sprintf("%v", pt)
// 	}
// 	return strings.Join(s, "->")
// }

// func (p *Path) Equals(m Move) bool {
// 	// check if m is a path
// 	p2, ok := m.(*Path)
// 	// check if p et p2 are the same length
// 	if ok && p.Len() != p2.Len() {
// 		ok = false
// 		fmt.Println("not the same length")
// 	}
// 	// check if each move is equal
// 	for i := 0; i < p.Len() && ok; i++ {
// 		ok = p.Moves[i].Equals(p2.Moves[i])
// 	}
// 	return ok
// }

func (p *Path) Start() Vector {
	return p.Moves[0].Start()
}

func (p *Path) End() Vector {
	return p.Moves[p.Len()-1].End()
}

// Reverse reverses path p, and all its composing moves
func (p *Path) Reverse() {
	for i, _ := range p.Moves {
		p.Moves[i].Reverse()
	}
	for i := 0; i < p.Len()/2; i++ {
		j := p.Len() - i - 1
		p.Moves[i], p.Moves[j] = p.Moves[j], p.Moves[i]
	}
}

// IsContinuous returns true if the path has no gaps, ie the next move always
// starts where the previous one ended
func (p *Path) IsContinuous() bool {
	ok := true
	for i := 0; i < p.Len()-1 && ok; i++ {
		ok = p.Moves[i].End() == p.Moves[i+1].Start()
	}
	return ok
}

// Points returns the list of points visited by the path as []Vector
func (p *Path) Points() []Vector {
	vs := make([]Vector, p.Len()+1)
	for i, m := range p.Moves {
		vs[i] = m.Start()
	}
	vs[p.Len()] = p.End()
	return vs
}

func (p *Path) Append(m Move) {
	p.Moves = append(p.Moves, m)
}

func (p *Path) Join(p2 *Path) {
	p.Moves = append(p.Moves, p2.Moves...)
}

// IsClosed returns true if p.Start() == p.End()
func (p *Path) IsClosed() bool {
	return p.Start() == p.End()
}

// IsClockwise returns true if the path is running clockwise, false otherwise.
// The shoelace algorithm is used to determine the direction of rotation
func (p *Path) IsClockwise() bool {
	sum := 0.0
	for _, m := range p.Moves {
		start := m.Start()
		end := m.End()
		sum += (end.X - start.X) * (end.Y + start.Y)
	}
	// the curve is CW if the sum is positive, CCW if the sum is negative
	return sum > 0
}

// HasInnerLoop checks if the path is passing through the same point twice.
func (p *Path) HasInnerLoop() (bool, int, int) {
	// a dictionary of visited starting points, and their index in the path
	visited := make(map[Vector]int)

	for i, move := range p.Moves {
		start := move.Start()
		if j, ok := visited[start]; !ok {
			// first time here ! store the starting point and corresponding index
			visited[start] = i
		} else {
			// found a loop ! starting from the stored index to the current index
			return true, j, i
		}
	}
	// special case : inner loop at the end of the path that doesn't start at
	// the begining
	if j, ok := visited[p.End()]; ok && j != 0 {
		return true, j, p.Len()
	}

	// no inner loop found
	return false, 0, 0
}

func (p *Path) SplitInnerLoops() []Path {
	split := []Path{}
	for p != nil && p.Len() > 0 {
		if ok, start, end := p.HasInnerLoop(); ok {
			// extract loop from p
			loop := &Path{make([]Move, end-start)}
			copy(loop.Moves, p.Moves[start:end])
			// merge remaining parts
			p.Moves = append(p.Moves[:start], p.Moves[end:]...)
			split = append(split, *loop)
		} else {
			split = append(split, *p)
			p = nil
		}
	}
	return split
}
