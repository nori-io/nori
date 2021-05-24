package plugin_option

import (
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/domain/service"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	PluginOptionRepository repository.PluginOptionRepository
}

func New(params Params) service.PluginOptionService {
	return &Service{
		PluginOptionRepository: params.PluginOptionRepository,
	}
}
