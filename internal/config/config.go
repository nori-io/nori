package config

import (
	go_config "github.com/cheebo/go-config"
	"github.com/cheebo/go-config/sources/file"
	"github.com/nori-io/logger"
)

type Config struct {
	Hooks   Hooks
	Plugins Plugins
	Storage Storage
}

type Hooks struct {
	Hooks []string
}

type Plugins struct {
	Dirs []string
}

type Storage struct {
	DSN string
}

func New(path string) func() (*Config, error) {
	return func() (*Config, error) {
		var c Config

		if path == "" {
			path = "/etc/nori/config.yml"
		}

		cfg := go_config.New()
		fileSource, err := file.Source(file.File{
			Path:      path,
			Namespace: "",
		})
		if err != nil {
			logger.L().Error(err.Error())
			return nil, err
		}
		cfg.UseSource(fileSource)

		err = cfg.Unmarshal(&c, "")
		if err != nil {
			logger.L().Error(err.Error())
			return nil, err
		}

		return &c, nil
	}
}
