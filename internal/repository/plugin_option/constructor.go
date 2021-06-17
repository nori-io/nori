package plugin_option

import (
	"github.com/nori-io/common/v5/pkg/domain/logger"
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/env"
	"go.uber.org/dig"
)

type Params struct {
	dig.In

	Env    *env.Env
	Logger logger.Logger
}

func New(params Params) (repository.PluginOptionRepository, error) {
	bucket, err := params.Env.Storage.CreateBucketIfNotExists("plugin_options")
	if err != nil {
		return nil, err
	}

	return &Repository{
		Bucket: bucket,
	}, nil
}
