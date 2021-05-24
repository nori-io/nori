package config

import (
	go_config "github.com/cheebo/go-config"
	"github.com/cheebo/go-config/sources/env"
	"github.com/cheebo/go-config/sources/file"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/registry"
)

type Params struct {
	File string
}

func New(params Params) (registry.ConfigRegistry, error) {
	config := go_config.New()
	config.UseSource(env.Source("NORI", "_"))

	if params.File != "" {
		fileSource, err := file.Source(file.File{
			Path: params.File,
		})
		if err != nil {
			return nil, err
		}
		config.UseSource(fileSource)
	}

	return &ConfigRegistry{
		pluginVariables: map[meta.ID]*[]registry.Variable{},
		config:          config,
	}, nil
}
