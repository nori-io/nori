package dependency_graph

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/registry"
	"github.com/nori-io/nori/pkg/nori/helper/dependency_graph/graph"
)

func New(registry registry.Registry) *DependencyGraphHelper {
	return &DependencyGraphHelper{
		graph:      graph.NewGraph(),
		registry:   registry,
		unresolved: map[meta.ID][]meta.Dependency{},
	}
}
