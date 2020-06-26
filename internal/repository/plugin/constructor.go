package plugin

import (
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/repository/plugin/graph"
)

func New() repository.PluginRepository {
	return &PluginRepository{
		graph: graph.NewDependencyGraph(),
	}
}
