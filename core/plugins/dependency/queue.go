package dependency

import (
	"sync"

	"gonum.org/v1/gonum/graph"
)

// NodeQueue the queue of Items
type NodeQueue struct {
	items []graph.Node
	lock  sync.RWMutex
}

// New creates a new NodeQueue
func NewNodeQueue() *NodeQueue {
	return &NodeQueue{
		items: []graph.Node{},
		lock:  sync.RWMutex{},
	}
}

func (s *NodeQueue) Enqueue(t graph.Node) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

func (s *NodeQueue) Dequeue() *graph.Node {
	s.lock.Lock()
	item := s.items[0]
	s.items = s.items[1:len(s.items)]
	s.lock.Unlock()
	return &item
}

func (s *NodeQueue) IsEmpty() bool {
	return len(s.items) == 0
}

func (s *NodeQueue) Size() int {
	return len(s.items)
}
