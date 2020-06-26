package handler

import (
	"github.com/nori-io/nori/internal/handler/grpc"
	"go.uber.org/fx"
)

func New() fx.Option {
	return fx.Provide(
		grpc.New,
	)
}
