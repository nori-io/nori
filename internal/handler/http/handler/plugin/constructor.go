package plugin

import (
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env *env.Env
	PluginManager service.PluginManager
}

type Handler struct {
	Env *env.Env
	PluginManager service.PluginManager
}

func New(params Params) *Handler {
	return &Handler{
		Env:           params.Env,
		PluginManager: params.PluginManager,
	}
}
