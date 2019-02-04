// Copyright © 2018 Nori info@nori.io
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation, either version 3
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package graph

import (
	"fmt"

	"github.com/nori-io/nori-common/meta"
)

// Directed Graph
type DependencyGraph interface {
	// Adds a node to the graph
	AddNode(id meta.ID) error

	// removes a node from the graph
	// and any edges attached to it
	RemoveNode(id meta.ID)

	// returns whether the node exists within the graph
	Has(id meta.ID) bool

	// return the node
	Node(id meta.ID) meta.ID

	// returns all nodes in the graph
	Nodes() []meta.ID

	// returns the edge from `from` to `to` if such an edge
	// exists and nil otherwise. The node v must be directly
	// reachable from `from`
	Edge(from, to meta.ID) Edge

	// returns all graph edges
	Edges() []Edge

	// returns a new Edge from the source to the destination node.
	NewEdge(from, to meta.ID) Edge

	// adds an edge from one node to another.
	SetEdge(e Edge)

	// removes the given edge, leaving the terminal nodes
	RemoveEdge(Edge)

	// returns whether an edge exists
	HasEdgeFromTo(u, v meta.ID) bool

	// returns all nodes that can be reached directly
	// from the given node.
	From(id meta.ID) []meta.ID

	// returns all nodes that can reach directly
	// to the given node.
	To(id meta.ID) []meta.ID

	// "topological" sorting
	Sort() ([]meta.ID, error)
}

func NewDependencyGraph() DependencyGraph {
	return &dependencyGraph{
		nodes: make(map[meta.ID]meta.ID),
		from:  make(map[meta.ID]map[meta.ID]Edge),
		to:    make(map[meta.ID]map[meta.ID]Edge),
	}
}

type dependencyGraph struct {
	nodes map[meta.ID]meta.ID
	from  map[meta.ID]map[meta.ID]Edge
	to    map[meta.ID]map[meta.ID]Edge
}

func (g *dependencyGraph) AddNode(id meta.ID) error {
	if _, exists := g.nodes[id]; exists {
		return fmt.Errorf("node already exists")
	}
	g.nodes[id] = id
	g.from[id] = make(map[meta.ID]Edge)
	g.to[id] = make(map[meta.ID]Edge)

	return nil
}

func (g *dependencyGraph) RemoveNode(id meta.ID) {
	if _, ok := g.nodes[id]; !ok {
		return
	}
	delete(g.nodes, id)

	for from := range g.from[id] {
		delete(g.to[from], id)
	}
	delete(g.from, id)

	for to := range g.to[id] {
		delete(g.from[to], id)
	}
	delete(g.to, id)
}

func (g *dependencyGraph) Has(id meta.ID) bool {
	_, ok := g.nodes[id]

	return ok
}

func (g *dependencyGraph) Node(id meta.ID) meta.ID {
	return g.nodes[id]
}

func (g *dependencyGraph) Nodes() []meta.ID {
	nodes := make([]meta.ID, len(g.nodes))
	i := 0
	for _, n := range g.nodes {
		nodes[i] = n
		i++
	}

	return nodes
}

func (g *dependencyGraph) From(id meta.ID) []meta.ID {
	if _, ok := g.from[id]; !ok {
		return nil
	}

	from := make([]meta.ID, len(g.from[id]))
	i := 0
	for item := range g.from[id] {
		from[i] = item
		i++
	}

	return from
}

func (g *dependencyGraph) Edge(from, to meta.ID) Edge {
	if _, ok := g.nodes[from]; !ok {
		return nil
	}
	if _, ok := g.nodes[from]; !ok {
		return nil
	}
	edge, ok := g.from[from][to]
	if !ok {
		return nil
	}
	return edge
}

func (g *dependencyGraph) Edges() []Edge {
	var edges []Edge
	for _, id := range g.nodes {
		for _, e := range g.from[id] {
			edges = append(edges, e)
		}
	}
	return edges
}

func (g *dependencyGraph) HasEdgeFromTo(from, to meta.ID) bool {
	if _, ok := g.nodes[from]; !ok {
		return false
	}
	if _, ok := g.nodes[to]; !ok {
		return false
	}
	if _, ok := g.from[from][to]; !ok {
		return false
	}
	return true
}

func (g *dependencyGraph) To(id meta.ID) []meta.ID {
	if _, ok := g.from[id]; !ok {
		return nil
	}

	to := make([]meta.ID, len(g.to[id]))
	i := 0
	for item := range g.to[id] {
		to[i] = item
		i++
	}

	return to
}

func (g *dependencyGraph) NewEdge(from, to meta.ID) Edge {
	return &edge{from: from, to: to}
}

func (g *dependencyGraph) SetEdge(e Edge) {
	var (
		from = e.From()
		to   = e.To()
	)

	if from == to {
		return
	}

	if !g.Has(from) {
		g.AddNode(from)
	}
	if !g.Has(to) {
		g.AddNode(to)
	}

	g.from[from][to] = e
	g.to[to][from] = e
}

func (g *dependencyGraph) RemoveEdge(e Edge) {
	from, to := e.From(), e.To()
	if _, ok := g.nodes[from]; !ok {
		return
	}
	if _, ok := g.nodes[to]; !ok {
		return
	}

	delete(g.from[from], to)
	delete(g.to[to], from)
}

func (g *dependencyGraph) Sort() ([]meta.ID, error) {
	var sorted []meta.ID
	queue := NewNodeQueue()

	tmpGraph := NewDependencyGraph()
	for _, n := range g.Nodes() {
		tmpGraph.AddNode(n)
	}
	for _, e := range g.Edges() {
		tmpGraph.SetEdge(tmpGraph.NewEdge(e.From(), e.To()))
	}

	for _, n := range tmpGraph.Nodes() {
		if len(tmpGraph.From(n)) == 0 {
			queue.Enqueue(n)
		}
	}

	for {
		if queue.IsEmpty() {
			break
		}
		n := queue.Dequeue()

		// add n to tail of L
		sorted = append(sorted, *n)
		// for each node m with an edge e from n to m do
		for _, m := range tmpGraph.To(*n) {
			e := tmpGraph.Edge(m, *n)
			// remove edge e from the graph
			if e != nil {
				tmpGraph.RemoveEdge(e)
			}
			// if m has no other incoming edges then insert m into S
			if len(tmpGraph.From(m)) == 0 {
				queue.Enqueue(m)
			}
		}
	}

	if len(tmpGraph.Edges()) > 0 {
		// @todo return cycle info
		return []meta.ID{}, fmt.Errorf("dependency cycle found")
	}

	return sorted, nil
}