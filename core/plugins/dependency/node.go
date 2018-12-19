package dependency

import "github.com/secure2work/nori/core/plugins/meta"

type Node interface {
	ID() int
}

type PNode interface {
	PID() meta.ID
}

func NewNode(id int, ID meta.ID) Node {
	return node{
		id:  id,
		pid: ID,
	}
}

type node struct {
	id  int
	pid meta.ID
}

func (d node) ID() int {
	return d.id
}

func (d node) PID() meta.ID {
	return d.pid
}
