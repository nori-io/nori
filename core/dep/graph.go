package dep

import (
	"fmt"

	"github.com/gonum/graph"
	"github.com/gonum/graph/simple"
)

func topologicalSort(g *simple.DirectedGraph) ([]*graph.Node, error) {
	var sorted []*graph.Node
	q := NewNodeQueue()

	g2 := simple.NewDirectedGraph(0, 0)
	for _, n := range g.Nodes() {
		g2.AddNode(n)
	}
	for _, e := range g.Edges() {
		g2.SetEdge(simple.Edge{
			F: e.From(), T: e.To(), W: e.Weight(),
		})
	}

	for _, n := range g2.Nodes() {
		if len(g2.From(n)) == 0 {
			q.Enqueue(n)
		}
	}

	for {
		if q.IsEmpty() {
			break
		}
		n := q.Dequeue()

		// add n to tail of L
		sorted = append(sorted, n)
		// for each node m with an edge e from n to m do
		for _, m := range g2.To(*n) {
			e := g2.Edge(m, *n)
			// remove edge e from the graph
			if e != nil {
				g2.RemoveEdge(e)
			}
			// if m has no other incoming edges then insert m into S
			if len(g2.From(m)) == 0 {
				q.Enqueue(m)
			}
		}
	}

	if len(g2.Edges()) > 0 {
		// @todo return cycle info
		return []*graph.Node{}, fmt.Errorf("dependency cycle found")
	}

	return sorted, nil
}
