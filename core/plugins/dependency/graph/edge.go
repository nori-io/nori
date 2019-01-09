package graph

import "github.com/secure2work/nori/core/plugins/meta"

type Edge interface {
	From() meta.ID
	To() meta.ID
}

type edge struct {
	from meta.ID
	to   meta.ID
}

func (e *edge) From() meta.ID {
	return e.from
}

func (e *edge) To() meta.ID {
	return e.to
}
