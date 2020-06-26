package nori

import (
	"github.com/nori-io/nori/internal/nori"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Provide(
		nori.New,
	)
}
