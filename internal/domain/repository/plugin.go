package repository

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/internal/domain/entity"
)

type PluginRepository interface {
	Create(file *entity.File) (*entity.Plugin, error)
	Delete(file *entity.File) error

	Find(id meta.ID) (*entity.Plugin, error)
	FindAll() []*entity.Plugin
	FindByIDs(ids []meta.ID) []*entity.Plugin
}
