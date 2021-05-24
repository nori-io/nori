package plugin

import (
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/internal/domain/repository"
)

type PluginRepository struct {
	plugins map[string]*entity.Plugin
	files   map[string]*entity.Plugin
}

func New() repository.PluginRepository {
	return &PluginRepository{
		plugins: map[string]*entity.Plugin{},
		files:   map[string]*entity.Plugin{},
	}
}
