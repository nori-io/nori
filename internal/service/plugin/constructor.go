package plugin

import (
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/domain/service"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	PluginRepository repository.PluginRepository
}

type Service struct {
	PluginRepository repository.PluginRepository
}

func New(params Params) service.PluginService {
	return &Service{
		PluginRepository: params.PluginRepository,
	}
}
