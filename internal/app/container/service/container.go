package service

import (
	"github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/service/registry"
	"github.com/nori-io/nori/internal/service/storage"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Provide(
		registry.New,
		storage.New,
		logger.L,
	)
}
