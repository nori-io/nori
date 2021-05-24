package repository

import (
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/repository/file"
	"github.com/nori-io/nori/internal/repository/plugin"
	"github.com/nori-io/nori/internal/repository/plugin_option"
	"go.uber.org/dig"
)

func Provide(container *dig.Container) {
	if err := container.Provide(
		file.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	if err := container.Provide(
		plugin.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	if err := container.Provide(
		plugin_option.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}
}
