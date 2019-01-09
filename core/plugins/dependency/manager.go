package dependency

import (
	"github.com/secure2work/nori/core/plugins/dependency/graph"
	"github.com/secure2work/nori/core/plugins/errors"
	"github.com/secure2work/nori/core/plugins/meta"
)

type Manager interface {
	// adds meta to dependency manager
	Add(m meta.Meta) error
	// returns whether the ID exists in dependency list
	Has(id meta.ID) bool
	// removes meta from dependency list
	Remove(id meta.ID)
	// returns ID of plugin that fits given dependency spec
	Resolve(dependency meta.Dependency) (meta.ID, error)

	// returns list of unresolved dependencies,
	// unresolved dependency: no plugins match dependency spec
	UnResolvedDependencies() map[meta.ID][]meta.Dependency

	// returns list of running plugin IDs
	Running() []meta.ID
	// return whether the plugin running or not
	IsRunning(id meta.ID) bool

	// marks the plugin as started
	Start(id meta.ID)
	// marks the plugin as stopped
	Stop(id meta.ID)

	// returns sorted consistent list of ID ready to start
	Sort() ([]meta.ID, error)
}

type manager struct {
	plugins    map[meta.ID]meta.Meta
	running    map[meta.ID]bool
	errors     []error
	graph      graph.DependencyGraph
	unresolved map[meta.ID][]meta.Dependency
}

func NewManager() Manager {
	return &manager{
		plugins:    map[meta.ID]meta.Meta{},
		running:    map[meta.ID]bool{},
		errors:     []error{},
		graph:      graph.NewDependencyGraph(),
		unresolved: map[meta.ID][]meta.Dependency{},
	}
}

func (m *manager) Add(mt meta.Meta) error {
	m.plugins[mt.Id()] = mt

	// add to graph
	err := m.graph.AddNode(mt.Id())
	if err != nil {
		return err
	}

	// build edges
	for _, dep := range mt.GetDependencies() {
		depId, err := m.Resolve(dep)
		if err != nil {
			if _, ok := m.unresolved[mt.Id()]; !ok {
				m.unresolved[mt.Id()] = []meta.Dependency{}
			}
			m.unresolved[mt.Id()] = append(m.unresolved[mt.Id()], dep)
			continue
		}
		m.graph.SetEdge(m.graph.NewEdge(mt.Id(), depId))
	}

	for unId, deps := range m.unresolved {
		if mt.Id() == unId {
			continue
		}
		for i, dep := range deps {
			depId, err := m.Resolve(dep)
			if err != nil {
				continue
			}
			m.graph.SetEdge(m.graph.NewEdge(unId, depId))
			m.unresolved[unId] = append(deps[:i], deps[i+1:]...)
		}
		if len(m.unresolved[unId]) == 0 {
			delete(m.unresolved, unId)
		}
	}

	return nil
}

func (m *manager) Has(id meta.ID) bool {
	_, ok := m.plugins[id]
	return ok
}

func (m *manager) Remove(id meta.ID) {
	delete(m.unresolved, id)
	delete(m.plugins, id)
}

func (m *manager) Resolve(dependency meta.Dependency) (meta.ID, error) {
	for id := range m.plugins {
		if id.ID != dependency.ID {
			continue
		}
		constraints, err := dependency.GetConstraint()
		if err != nil {
			return meta.ID{}, err
		}
		version, err := id.GetVersion()
		if err != nil {
			return meta.ID{}, err
		}
		if constraints.Check(version) {
			return id, nil
		}
	}
	return meta.ID{}, errors.DependencyNotFound{
		Dependency: dependency,
	}
}

func (m *manager) UnResolvedDependencies() map[meta.ID][]meta.Dependency {
	unresolved := map[meta.ID][]meta.Dependency{}
	for id, deps := range m.unresolved {
		list := make([]meta.Dependency, len(deps))
		copy(list, deps)
		unresolved[id] = list
	}
	return unresolved
}

func (m *manager) Running() []meta.ID {
	var list []meta.ID
	for id := range m.running {
		list = append(list, id)
	}
	return list
}

func (m *manager) IsRunning(id meta.ID) bool {
	_, ok := m.running[id]
	return ok
}

func (m *manager) Start(id meta.ID) {
	m.running[id] = true
}

func (m *manager) Stop(id meta.ID) {
	delete(m.running, id)
}

func (m *manager) Sort() ([]meta.ID, error) {
	return m.graph.Sort()
}
