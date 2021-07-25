package dependency_graph

import (
	"errors"

	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/registry"
	"github.com/nori-io/nori/pkg/nori/helper/dependency_graph/graph"
)

type DependencyGraphHelper struct {
	registry   registry.Registry
	graph      graph.DirectedGraph
	unresolved map[meta.ID][]meta.Dependency
}

func (h *DependencyGraphHelper) Add(id meta.ID, deps []meta.Dependency) error {
	// check already registered
	if h.graph.Has(id) {
		return nil
	}

	// add to dependency graph
	if err := h.graph.AddNode(id); err != nil {
		return err
	}

	// build dependency graph edges for plugin
	unresolved, err := h.buildEdges(id, deps)
	if err != nil {
		return err
	}
	//if _, ok := h.unresolved[p.Meta().GetID()]; !ok {
	//	h.unresolved[p.Meta().GetID()] = []meta.Dependency{}
	//}
	if len(unresolved) > 0 {
		h.unresolved[id] = unresolved
	}

	// try to resolve and build dependency edges for unresolved dependencies
	return h.resolveUnresolvedDeps(id)
}

func (h *DependencyGraphHelper) Remove(id meta.ID) error {
	delete(h.unresolved, id)
	if !h.graph.Has(id) {
		return nil
	}
	h.graph.RemoveNode(id)
	return nil
}

func (h *DependencyGraphHelper) Has(id meta.ID) bool {
	return h.graph.Has(id)
}

func (h *DependencyGraphHelper) HasUnresolved() bool {
	return len(h.unresolved) > 0
}

func (h *DependencyGraphHelper) GetSorted() ([]meta.ID, error) {
	return h.graph.Sort()
}

func (h *DependencyGraphHelper) GetSortedReversed() ([]meta.ID, error) {
	list, err := h.graph.Sort()
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list, nil
}

func (h *DependencyGraphHelper) GetUnresolved() map[meta.ID][]meta.Dependency {
	return h.unresolved
}

func (h *DependencyGraphHelper) GetDependencies(id meta.ID) []meta.ID {
	return h.graph.From(id)
}

func (h *DependencyGraphHelper) GetDependent(id meta.ID) []meta.ID {
	return h.graph.To(id)
}

func (h *DependencyGraphHelper) buildEdges(id meta.ID, deps []meta.Dependency) ([]meta.Dependency, error) {
	unresolved := []meta.Dependency{}
	for _, dep := range deps {
		resolves := h.registry.ResolveDependency(id, dep)
		if len(resolves) > 1 {
			// todo: not supported yet
			return nil, errors.New("multiple resolves found for dependency " + dep.String())
		}
		if len(resolves) == 0 {
			unresolved = append(unresolved, dep)
			continue
		}
		h.graph.SetEdge(h.graph.NewEdge(id, resolves[0].Meta().GetID()))
	}
	return unresolved, nil
}

func (h *DependencyGraphHelper) resolveUnresolvedDeps(pid meta.ID) error {
	for id, deps := range h.unresolved {
		if pid.String() == id.String() {
			continue
		}
		unresolved, err := h.buildEdges(id, deps)
		if err != nil {
			return err
		}
		if len(unresolved) > 0 {
			h.unresolved[id] = unresolved
		}
		if len(unresolved) == 0 {
			delete(h.unresolved, id)
		}
	}

	return nil
}
