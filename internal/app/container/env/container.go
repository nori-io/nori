package env

import (
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/config"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/fx"
)

func New(path string, log logger.Logger) func() fx.Option {
	return func() fx.Option {
		return fx.Provide(
			env.New,
			config.New(path),
		)
	}
}
