package handler

import (
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/handler/grpc"
	"github.com/nori-io/nori/internal/handler/http"
	"github.com/nori-io/nori/internal/handler/http/handler/plugin"
	"go.uber.org/dig"
)

func Provide(container *dig.Container) {
	if err := container.Provide(
		grpc.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	if err := container.Provide(
		http.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	// http handlers
	if err := container.Provide(
		plugin.New,
	); err != nil {
		log.L().Fatal("%s", err.Error())
	}
}
