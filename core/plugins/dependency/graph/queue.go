package graph

import (
	"sync"

	"github.com/secure2work/nori/core/plugins/meta"
)

// NodeQueue the queue of Items
type NodeQueue struct {
	items []meta.ID
	lock  sync.RWMutex
}

// New creates a new NodeQueue
func NewNodeQueue() *NodeQueue {
	return &NodeQueue{
		items: []meta.ID{},
		lock:  sync.RWMutex{},
	}
}

func (s *NodeQueue) Enqueue(t meta.ID) {
	s.lock.Lock()
	s.items = append(s.items, t)
	s.lock.Unlock()
}

func (s *NodeQueue) Dequeue() *meta.ID {
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
