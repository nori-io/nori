// Copyright © 2018 Nori info@nori.io
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation, either version 3
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package dependency

import (
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori/internal/dependency/graph"
	"github.com/nori-io/nori/pkg/errors"
	"github.com/nori-io/nori/pkg/types"
)

type Manager interface {
	// adds meta to dependency manager
	Add(m meta.Meta) error

	GetPluginsList() map[meta.ID]meta.Meta

	// get dependent plugins ids
	GetDependent(id meta.ID) []meta.ID

	// returns whether the ID exists in dependency list
	Has(id meta.ID) bool

	// removes meta from dependency list
	Remove(id meta.ID)

	// returns ID of plugin that fits given dependency spec
	Resolve(dependency meta.Dependency) (meta.ID, error)

	// returns list of unresolved dependencies,
	// unresolved dependency: no plugins match dependency spec
	UnResolvedDependencies() map[meta.ID][]meta.Dependency

	// returns sorted consistent list of ID ready to start
	Sort() ([]meta.ID, error)
}

type manager struct {
	plugins    map[meta.ID]meta.Meta
	running    types.MetaList
	errors     []error
	graph      graph.DependencyGraph
	unresolved map[meta.ID][]meta.Dependency
}

func NewManager() Manager {
	return &manager{
		plugins:    map[meta.ID]meta.Meta{},
		running:    types.MetaList{},
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
		depID, err := m.Resolve(dep)
		if err != nil {
			if _, ok := m.unresolved[mt.Id()]; !ok {
				m.unresolved[mt.Id()] = []meta.Dependency{}
			}
			m.unresolved[mt.Id()] = append(m.unresolved[mt.Id()], dep)
			continue
		}
		if depID.ID == mt.Id().ID {
			m.unresolved[mt.Id()] = append(m.unresolved[mt.Id()], dep)
			continue
		}
		m.graph.SetEdge(m.graph.NewEdge(mt.Id(), depID))
	}

	for unID, deps := range m.unresolved {
		if mt.Id() == unID {
			continue
		}
		for i, dep := range deps {
			depID, err := m.Resolve(dep)

			if err != nil {
				continue
			}
			m.graph.SetEdge(m.graph.NewEdge(unID, depID))
			m.unresolved[unID] = append(deps[:i], deps[i+1:]...)
		}
		if len(m.unresolved[unID]) == 0 {
			delete(m.unresolved, unID)
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
	// @todo delete node and related edges?
	m.graph.RemoveNode(id)
}

func (m *manager) Resolve(dependency meta.Dependency) (meta.ID, error) {
	for id, m := range m.plugins {
		// dependency on interface
		// dependency on plugin
		if id.ID != dependency.ID {
			if m.GetInterface() != dependency.Interface {
				continue
			}
		}

		if !dependency.Interface.IsUndefined() && m.GetInterface().Equal(dependency.Interface) {
			return id, nil
		}

		if id.ID == dependency.ID {
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

func (m *manager) Sort() ([]meta.ID, error) {
	if len(m.unresolved) > 0 {
		return nil, errors.DependenciesNotFound{
			Dependencies: m.unresolved,
		}
	}
	return m.graph.Sort()
}

func (m *manager) GetPluginsList() map[meta.ID]meta.Meta {
	return m.plugins
}

func (m *manager) GetDependent(id meta.ID) []meta.ID {
	if m.graph.Has(id) {
		return m.graph.To(id)
	}
	return []meta.ID{}
}
