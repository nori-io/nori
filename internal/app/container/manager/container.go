package manager

import (
	"github.com/nori-io/nori/internal/manager/config"
	"github.com/nori-io/nori/internal/manager/file"
	"github.com/nori-io/nori/internal/manager/plugin"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Provide(
		file.New,
		plugin.New,
		config.New,
	)
}
