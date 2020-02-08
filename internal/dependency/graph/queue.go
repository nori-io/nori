/*
Copyright 2019-2020 The Nori Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graph

import (
	"sync"

	"github.com/nori-io/nori-common/v2/meta"
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
