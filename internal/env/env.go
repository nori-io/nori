package env

import (
	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/common/v5/pkg/domain/storage"
	"github.com/nori-io/nori/internal/config"
	store "github.com/nori-io/nori/internal/env/storage"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Config *config.Config
	Logger logger.Logger
}

type Env struct {
	Config  *config.Config
	Logger  logger.Logger
	Storage storage.Storage
}

func New(params Params) (*Env, error) {
	store, err := store.New(params.Config)
	if err != nil {
		return nil, err
	}

	return &Env{
		Config:  params.Config,
		Logger:  params.Logger,
		Storage: store,
	}, nil
}
