package repository

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/internal/domain/entity"
)

type PluginOptionRepository interface {
	Upsert(po entity.PluginOption) error
	Delete(id meta.ID) error

	Find(id meta.ID) (entity.PluginOption, error)
	FindAll() ([]entity.PluginOption, error)
}
