package service

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/internal/domain/entity"
)

type PluginOptionService interface {
	Upsert(data PluginOptionUpsertData) (entity.PluginOption, error)
	Delete(id meta.ID) error

	Get(id meta.ID) (entity.PluginOption, error)
	GetAll() ([]entity.PluginOption, error)
}

type PluginOptionUpsertData struct {
	ID          meta.ID
	Enabled     bool
	Installed   bool
	Installable bool
}
