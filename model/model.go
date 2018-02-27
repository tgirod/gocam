package model

type Model struct {
	Paths []Path
}

func New() *Model {
	return new(Model)
}

func (m *Model) Add(p *Path) {
	for p != nil {
		found := false
		for i := 0; i < len(m.Paths) && !found; i++ {
			q := m.Paths[i]
			if !q.IsClosed() {
				if q.Start() == p.End() {
					m.Paths = append(m.Paths[:i], m.Paths[i+1:]...) // remove q from paths
					p.Join(&q)                                      // join paths
					found = true                                    // exit loop
				} else if q.Start() == p.Start() {
					p.Reverse()
					m.Paths = append(m.Paths[:i], m.Paths[i+1:]...) // remove q from paths
					p.Join(&q)                                      // join paths
					found = true                                    // exit loop
				}
			}
		}
		if !found {
			// no matching path - add p to the paths and mark it nil to leave the loop
			m.Paths = append(m.Paths, *p)
			p = nil
		}
	}
}
