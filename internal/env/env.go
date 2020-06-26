package env

import (
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori/internal/config"
	"go.uber.org/fx"
)

type Env struct {
	Config *config.Config
	Logger logger.Logger
}

type Params struct {
	fx.In

	Config *config.Config
}

func New(params Params) (*Env, error) {
	return &Env{
		Config: params.Config,
		Logger: log.L(),
	}, nil
}
