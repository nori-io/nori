package nori

import (
	log "github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/storage"
	"github.com/nori-io/nori/internal/domain/manager"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/fx"
)

type NoriParams struct {
	fx.In

	Env           *env.Env
	Logger        log.Logger
	FileManager   manager.File
	PluginManager manager.Plugin
	Storage       storage.Storage
}

func New(params NoriParams) (*Nori, error) {
	return &Nori{
		env: params.Env,
		log: params.Logger,
		managers: struct {
			File   manager.File
			Plugin manager.Plugin
		}{
			File:   params.FileManager,
			Plugin: params.PluginManager,
		},
		storage: params.Storage,
	}, nil
}
