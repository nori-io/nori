package service

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
)

type PluginService interface {
	Add(plugin *entity.Plugin) error

	Get(id meta.ID) (*entity.Plugin, error)
	GetAll() []*entity.Plugin
	GetByIDs(ids []meta.ID) []*entity.Plugin
}
