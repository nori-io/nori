package repository

import (
	"github.com/nori-io/nori/internal/repository/file"
	"github.com/nori-io/nori/internal/repository/plugin"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Provide(
		plugin.New,
		file.New,
	)
}
