package config

import (
	"os"

	go_config "github.com/cheebo/go-config"
	"github.com/cheebo/go-config/pkg/sources/env"
	"github.com/cheebo/go-config/pkg/sources/file"
	"github.com/nori-io/logger"
)

const (
	defaultConfig = "/etc/nori/config.yml"
)

type Config struct {
	Nori    Nori
	Hooks   go_config.Fields
	Plugins go_config.Fields
}

type Nori struct {
	Hooks   Hooks
	Plugins Plugins
	GRPC    GRPC
	Storage Storage
}

type Hooks struct {
	Dir string
}

type Plugins struct {
	Dir string
}

type GRPC struct {
	Host string
	Port uint
}

type Storage struct {
	DSN string
}

func New(path string) (*Config, error) {
	config := go_config.New()

	if path == "" {
		if _, err := os.Stat(defaultConfig); err != nil {
			return nil, err
		}
		path = defaultConfig
	}

	if path != "" {
		fileSource, err := file.Source(file.File{
			Path:      path,
			Namespace: "",
		})
		if err != nil {
			return nil, err
		}
		config.UseSource(fileSource)
	}

	config.UseSource(env.Source("NORI", "_"))

	config.SetDefault("nori.grpc.host", "0.0.0.0")
	config.SetDefault("nori.grpc.port", 5876)

	var n Nori
	if err := config.Sub("nori").Unmarshal(&n, ""); err != nil {
		logger.L().Fatal(err.Error())
		return nil, err
	}

	return &Config{
		Nori:    n,
		Hooks:   config.Sub("hooks"),
		Plugins: config.Sub("plugins"),
	}, nil
}
