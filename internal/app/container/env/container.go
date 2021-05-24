package env

import (
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/config"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/dig"
)

func Provide(container *dig.Container, configFile string) {
	// config
	if err := container.Provide(
		func() (*config.Config, error) { return config.New(configFile) },
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	// env
	if err := container.Provide(
		env.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}
}
