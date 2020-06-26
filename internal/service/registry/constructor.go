package registry

import (
	"github.com/nori-io/nori-common/v2/plugin"
	"github.com/nori-io/nori/internal/domain/repository"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	PluginRepository repository.PluginRepository
}

func New(params Params) (plugin.Registry, error) {
	return &Service{
		PluginRepository: params.PluginRepository,
	}, nil
}
