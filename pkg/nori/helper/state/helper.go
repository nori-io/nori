package state

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/errors"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
)

type StateHelper struct {
	plugins map[meta.ID]enum.State
}

func (s *StateHelper) GetState(id meta.ID) (enum.State, error) {
	status, ok := s.plugins[id]
	if !ok {
		return enum.Undefined, errors.NotFound{ID: id}
	}
	return status, nil
}

func (s *StateHelper) SetState(id meta.ID, state enum.State) {
	if s.plugins == nil {
		s.plugins = map[meta.ID]enum.State{}
	}
	s.plugins[id] = state
}

func (s *StateHelper) GetAllByState(state enum.State) []meta.ID {
	ids := []meta.ID{}
	for i, s := range s.plugins {
		if s == state {
			ids = append(ids, i)
		}
	}
	return ids
}
