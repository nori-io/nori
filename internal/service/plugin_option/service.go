package plugin_option

import (
	"time"

	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/errors"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/domain/service"
)

type Service struct {
	PluginOptionRepository repository.PluginOptionRepository
}

func (s Service) Upsert(data service.PluginOptionUpsertData) (entity.PluginOption, error) {
	po, err := s.PluginOptionRepository.Find(data.ID)
	if err != nil {
		if _, ok := err.(errors.EntityNotFound); !ok {
			return entity.PluginOption{}, err
		}
	}

	created := time.Now()
	if !po.IsEmpty() {
		created = po.CreatedAt
	}

	po = entity.PluginOption{
		ID:          data.ID,
		Enabled:     data.Enabled,
		Installed:   data.Installed,
		Installable: data.Installable,
		CreatedAt:   created,
		UpdatedAt:   time.Now(),
	}

	err = s.PluginOptionRepository.Upsert(po)
	return po, err
}

func (s Service) Delete(id meta.ID) error {
	return s.PluginOptionRepository.Delete(id)
}

func (s Service) Get(id meta.ID) (entity.PluginOption, error) {
	return s.PluginOptionRepository.Find(id)
}

func (s Service) GetAll() ([]entity.PluginOption, error) {
	return s.PluginOptionRepository.FindAll()
}
