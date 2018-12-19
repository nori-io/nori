package dependency

import "github.com/secure2work/nori/core/plugins/meta"

type Node interface {
	ID() int64
}

type PNode interface {
	PID() meta.ID
}

func NewNode(id int64, ID meta.ID) Node {
	return node{
		id:  id,
		pid: ID,
	}
}

type node struct {
	id  int64
	pid meta.ID
}

func (d node) ID() int64 {
	return d.id
}

func (d node) PID() meta.ID {
	return d.pid
}
