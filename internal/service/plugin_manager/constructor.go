package plugin_manager

import (
	"context"

	"github.com/nori-io/common/v5/pkg/domain/plugin"
	errors2 "github.com/nori-io/common/v5/pkg/errors"
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/domain/service"
	"github.com/nori-io/nori/internal/env"
	"github.com/nori-io/nori/pkg/nori"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env                 *env.Env
	FileService         service.FileService
	PluginService       service.PluginService
	PluginOptionService service.PluginOptionService
}

type PluginManager struct {
	Env                 *env.Env
	FileService         service.FileService
	PluginService       service.PluginService
	PluginOptionService service.PluginOptionService
	Nori                nori.Nori
}

func New(params Params) (service.PluginManager, error) {
	files, err := params.FileService.GetAll(params.Env.Config.Nori.Plugins.Dir)
	if err != nil {
		return nil, err
	}

	n, err := nori.New(log.L())
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p, err := params.PluginService.Create(&file)
		if err != nil {
			params.Env.Logger.Error(err.Error())
		}

		params.Env.Logger.Info("found %s (%s) in %s", p.Plugin.Meta().GetID().String(), p.Plugin.Meta().GetInterface().String(), file.Path)

		pluginOptions, err := params.PluginOptionService.Get(p.Plugin.Meta().GetID())

		if err != nil {
			// 'not found' is equal 'not enabled'
			if _, ok := err.(errors2.EntityNotFound); ok {
				continue
			}
			return nil, err
		}

		if _, ok := p.Plugin.(plugin.Installable); ok {
			// not installed
			if !pluginOptions.Installed {
				continue
			}
		}

		// not enabled
		if !pluginOptions.Enabled {
			continue
		}

		if err := n.Add(p.Plugin); err != nil {
			params.Env.Logger.Error("%s", err.Error())
		}
		params.Env.Logger.Info("loaded %s", file.Path)
	}

	pluginManager := &PluginManager{
		Env:                 params.Env,
		FileService:         params.FileService,
		PluginService:       params.PluginService,
		PluginOptionService: params.PluginOptionService,
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
