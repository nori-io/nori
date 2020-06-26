package plugin

import (
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori-common/v2/version"
	"github.com/nori-io/nori/internal/domain/entity"
	"github.com/nori-io/nori/internal/repository/plugin/graph"
	"github.com/nori-io/nori/pkg/errors"
)

type PluginRepository struct {
	plugins    []*entity.Plugin
	graph      graph.DependencyGraph
	unresolved map[meta.ID][]meta.Dependency
}

func (r *PluginRepository) Register(p *entity.Plugin) error {
	item := r.FindByID(p.Meta().Id())

	// plugin already in the list
	if item != nil {
		// todo: return error or log
		return nil
	}
	r.plugins = append(r.plugins, p)

	if err := r.registerDependency(p); err != nil {
		r.UnRegister(p) // ?
	}

	return nil
}

func (r *PluginRepository) UnRegister(p *entity.Plugin) error {
	for i, item := range r.plugins {
		if item.Meta().Id() == p.Meta().Id() {
			r.plugins = append(r.plugins[:i], r.plugins[:i+1]...)
			break
		}
	}

	return nil
}

func (r *PluginRepository) FindAll() []*entity.Plugin {
	return r.plugins
}

func (r *PluginRepository) FindByID(id meta.ID) *entity.Plugin {
	for _, p := range r.plugins {
		if p.Meta().Id().String() == id.String() {
			return p
		}
	}
	return nil
}

func (r *PluginRepository) FindByIDs(ids []meta.ID) []*entity.Plugin {
	items := []*entity.Plugin{}
	for _, p := range r.plugins {
		for _, id := range ids {
			if p.Meta().Id().String() == id.String() {
				items = append(items, p)
			}
		}
	}
	return items
}

func (r *PluginRepository) FindByInterface(i meta.Interface) []*entity.Plugin {
	var plugins []*entity.Plugin
	for _, p := range r.plugins {
		if p.Meta().GetInterface().Equal(i) {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// Resolve returns plugin that fits given dependency spec
func (r *PluginRepository) Resolve(dependency meta.Dependency) (*entity.Plugin, error) {
	for _, p := range r.plugins {
		// interface is undefined
		if dependency.Interface.IsUndefined() {
			continue
		}

		// names of interfaces are not equal
		if p.Meta().GetInterface().Name() != dependency.Interface.Name() {
			continue
		}

		constraint, err := dependency.GetConstraint()
		if err != nil {
			return nil, err
		}
		ver, err := version.NewVersion(p.Meta().GetInterface().Version())
		if err != nil {
			return nil, err
		}
		if constraint.Check(ver) {
			return p, nil
		}
	}

	return nil, errors.DependencyNotFound{
		Dependency: dependency,
	}
}

//
func (r *PluginRepository) FindDependent(id meta.ID) []*entity.Plugin {
	return r.FindByIDs(r.graph.To(id))
}

func (r *PluginRepository) FindInstallable() []*entity.Plugin {
	var plugins []*entity.Plugin
	for _, p := range r.plugins {
		if p.IsInstallable() {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

func (r *PluginRepository) FindHooks() []*entity.Plugin {
	var plugins []*entity.Plugin
	for _, p := range r.plugins {
		if p.IsHook() {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// Tree returns plugin execution tree as an array
func (r *PluginRepository) Tree() ([]*entity.Plugin, error) {
	if len(r.unresolved) > 0 {
		return nil, errors.DependenciesNotFound{
			Dependencies: r.unresolved,
		}
	}
	ids, err := r.graph.Sort()
	if err != nil {
		return nil, err
	}

	// keep ids and plugins order
	items := []*entity.Plugin{}
	for _, id := range ids {
		for _, p := range r.plugins {
			if p.Meta().Id().String() == id.String() {
				items = append(items, p)
			}
		}
	}
	return items, nil
}

func (r *PluginRepository) registerDependency(p *entity.Plugin) error {
	// dependency loop: self-dependency
	//for _, dep := range p.Meta().GetDependencies() {
	//	if p.Meta().GetInterface() == dep.Interface {
	//		loopVertex := dep
	//		return errors.LoopVertexFound{Dependency: loopVertex}
	//	}
	//}

	// add to dependency graph
	err := r.graph.AddNode(p.Meta().Id())
	if err != nil {
		return err
	}
	// build dependency graph edges
	for _, dep := range p.Meta().GetDependencies() {
		d, err := r.Resolve(dep)
		if err != nil {
			if _, ok := r.unresolved[p.Meta().Id()]; !ok {
				r.unresolved[p.Meta().Id()] = []meta.Dependency{}
			}
			r.unresolved[p.Meta().Id()] = append(r.unresolved[p.Meta().Id()], dep)
			continue
		}
		if d.Meta().Id().String() == p.Meta().Id().String() {
			r.unresolved[p.Meta().Id()] = append(r.unresolved[p.Meta().Id()], dep)
			continue
		}
		r.graph.SetEdge(r.graph.NewEdge(p.Meta().Id(), d.Meta().Id()))
	}
	// check unresolved dependencies
	for id, deps := range r.unresolved {
		if p.Meta().Id().String() == id.String() {
			continue
		}
		for i, dep := range deps {
			d, err := r.Resolve(dep)
			if err != nil {
				continue
			}

			r.graph.SetEdge(r.graph.NewEdge(id, d.Meta().Id()))
			r.unresolved[id] = append(deps[:i], deps[i+1:]...)
		}
		if len(r.unresolved[id]) == 0 {
			delete(r.unresolved, id)
		}
	}

	return nil
}

func (r *PluginRepository) unRegisterDependency(p *entity.Plugin) error {
	// remove dependencies
	delete(r.unresolved, p.Meta().Id())
	// @todo delete node and related edges?
	r.graph.RemoveNode(p.Meta().Id())

	// todo: check unresolved dependencies on dependent plugins

	return nil
}
