package plugin

import (
	"github.com/nori-io/nori-common/v2/storage"
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/repository/plugin/graph"
	"go.uber.org/fx"
)

type Params struct {
	fx.In
	Storage storage.Storage
}

func New(params Params) (repository.PluginRepository, error) {
	bucket, err := params.Storage.CreateBucketIfNotExists("plugins")
	if err != nil {
		return nil, err
	}

	return &PluginRepository{
		graph:  graph.NewDependencyGraph(),
		bucket: bucket,
	}, nil
}
