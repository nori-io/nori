package service

import (
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/service/file"
	"github.com/nori-io/nori/internal/service/hooks"
	"github.com/nori-io/nori/internal/service/plugin"
	"github.com/nori-io/nori/internal/service/plugin_manager"
	"github.com/nori-io/nori/internal/service/plugin_option"
	"go.uber.org/dig"
)

func Provide(container *dig.Container) {
	if err := container.Provide(
		file.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	if err := container.Provide(
		hooks.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	if err := container.Provide(
		plugin.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	if err := container.Provide(
		plugin_manager.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	if err := container.Provide(
		plugin_option.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}
}
