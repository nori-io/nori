package container

import (
	"github.com/nori-io/common/v5/pkg/domain/logger"
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/app"
	"github.com/nori-io/nori/internal/app/container/env"
	"github.com/nori-io/nori/internal/app/container/handler"
	"github.com/nori-io/nori/internal/app/container/repository"
	"github.com/nori-io/nori/internal/app/container/service"
	"go.uber.org/dig"
)

func New(configFile string) *dig.Container {
	container := dig.New()

	// FieldLogger
	if err := container.Provide(func() logger.Logger { return log.L() }); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	// app
	if err := container.Provide(app.New); err != nil {
		log.L().Fatal("%s", err.Error())
	}

	// config & env
	env.Provide(container, configFile)

	// repos
	repository.Provide(container)

	// services
	service.Provide(container)

	// handler
	handler.Provide(container)

	return container
}
