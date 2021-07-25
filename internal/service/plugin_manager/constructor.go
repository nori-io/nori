package plugin_manager

import (
	"context"
	"errors"

	"github.com/nori-io/common/v5/pkg/domain/registry"
	common_errors "github.com/nori-io/common/v5/pkg/errors"
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	"github.com/nori-io/nori/pkg/nori"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env                 *env.Env
	FileService         service.FileService
	PluginService       service.PluginService
	PluginOptionService service.PluginOptionService
	ConfigRegistry registry.ConfigRegistry
}

type PluginManager struct {
	Env                 *env.Env
	FileService         service.FileService
	PluginService       service.PluginService
	PluginOptionService service.PluginOptionService
	ConfigRegistry registry.ConfigRegistry
	Nori                nori.Nori
}

func New(params Params) (service.PluginManager, error) {
	files, err := params.FileService.GetAll(params.Env.Config.Nori.Plugins.Dir)

	if err != nil {
		params.Env.Logger.Fatal(err.Error())
	}

	n, err := nori.New(params.ConfigRegistry, log.L())
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p, err := entity.New(entity.File{
			Path: file.Path,
			Fn:   file.Fn,
		})
		if err != nil {
			params.Env.Logger.Error(err.Error())
		}

		params.Env.Logger.Info("found %s (%s) in %s", p.Meta().GetID().String(), p.Meta().GetInterface().String(), file.Path)

		pluginOptions, err := params.PluginOptionService.Get(p.Meta().GetID())
		if errors.Is(err, common_errors.EntityNotFound{}) {
				continue
		}
		if err != nil {
			return nil, err
		}

		if p.IsInstallable() {
			// not installed
			if !pluginOptions.Installed {
				continue
			}
		}

		// not enabled
		if !pluginOptions.Enabled {
			continue
		}

		if err := n.Add(p); err != nil {
			params.Env.Logger.Error("%s", err.Error())
		}
		params.Env.Logger.Info("loaded %s", file.Path)
	}

	pluginManager := &PluginManager{
		Env:                 params.Env,
		FileService:         params.FileService,
		PluginService:       params.PluginService,
		PluginOptionService: params.PluginOptionService,
		ConfigRegistry: params.ConfigRegistry,
		Nori:                n,
	}

	// todo: read action from config
	ctx := context.Background()
	err = n.InitAll(ctx)
	if err != nil {
		params.Env.Logger.Error(err.Error())
	} else {
		if err := n.StartAll(ctx); err != nil {
			params.Env.Logger.Error(err.Error())
		}
	}

	return pluginManager, nil
}
