package plugin

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/pkg/errors"
)

func (s *Service) Create(file *entity.File) (*entity.Plugin, error) {
	return s.PluginRepository.Create(file)
}

func (s *Service) Get(id meta.ID) (*entity.Plugin, error) {
	p, err := s.PluginRepository.Find(id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.NotFound{ID: id}
	}
	return p, nil
}

func (s *Service) GetAll() []*entity.Plugin {
	return s.PluginRepository.FindAll()
}

func (s *Service) GetByIDs(ids []meta.ID) []*entity.Plugin {
	return s.PluginRepository.FindByIDs(ids)
}
