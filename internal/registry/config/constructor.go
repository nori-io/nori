package config

import (
	go_config "github.com/cheebo/go-config"
	"github.com/cheebo/go-config/pkg/sources/env"
	"github.com/cheebo/go-config/pkg/sources/file"
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/registry"
)

func New(configFile string) (registry.ConfigRegistry, error) {
	config := go_config.New()
	config.UseSource(env.Source("NORI", "_"))

	if configFile != "" {
		fileSource, err := file.Source(file.File{
			Path: configFile,
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
