package plugin

import (
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env *env.Env
	FileService      service.FileService
	InstalledService service.PluginOptionService
	PluginManager    service.PluginManager
	PluginService    service.PluginService
}

type Handler struct {
	Env *env.Env
	FileService      service.FileService
	InstalledService service.PluginOptionService
	PluginManager    service.PluginManager
	PluginService    service.PluginService
}

func New(params Params) Handler {
	return Handler{
		Env:           params.Env,
		FileService: params.FileService,
		InstalledService: params.InstalledService,
		PluginManager: params.PluginManager,
		PluginService: params.PluginService,
	}
}
