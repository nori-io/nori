// Copyright © 2018 Secure2Work info@secure2work.com
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
	"sync"

	"github.com/secure2work/nori-common/meta"
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
