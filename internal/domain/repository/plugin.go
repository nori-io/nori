package repository

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
)

type PluginRepository interface {
	Add(plugin *entity.Plugin) error
	Remove(file string) error

	Find(id meta.ID) (*entity.Plugin, error)
	FindAll() []*entity.Plugin
	FindByIDs(ids []meta.ID) []*entity.Plugin
}
