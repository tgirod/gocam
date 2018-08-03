package util

type Graph struct {
	Nodes map[Node]Connect
	Edges map[Edge]bool
}

type Connect struct {
	Out []Edge
	In  []Edge
}

type Node interface{}

type Edge interface {
	From() Node
	To() Node
}

// AddNode adds node n to the graph if it does not already exists
func (g *Graph) AddNode(n Node) {
	if _, found := g.Nodes[n]; !found {
		g.Nodes[n] = Connect{}
	}
}

func (g *Graph) AddNodes(ns ...Node) {
	for _, n := range ns {
		g.AddNode(n)
	}
}

// AddEdge adds edge e to the graph if it does not already exists. It also adds
// From and To nodes if necessary
func (g *Graph) AddEdge(e Edge) {
	if _, found := g.Edges[e]; !found {
		g.Edges[e] = true   // add edge
		g.AddNode(e.From()) // optionally add e.From
		g.AddNode(e.To())   // optionally add e.To
		from := g.Nodes[e.From()]
		to := g.Nodes[e.To()]
		from.Out = append(from.Out, e) // connect to the From Node
		to.In = append(to.In, e)       // connect to the End Node
	}
}

func (g *Graph) DFS(start Node) {
	stack := []Node{start}
	visited := make(map[Node]Edge)
	for len(stack) > 0 {
		n := stack[0]
		stack = stack[1:]
		if _, found := visited[n]; !found {
			// never visited
		}
	}
}

// FindCycles runs a depth first search on the graph in order to find every
// possible cycle
func (g *Graph) FindCycles() [][]Edge {
	return nil // FIXME
}
