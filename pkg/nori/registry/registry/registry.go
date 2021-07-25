package registry

import (
	"fmt"
	"sync"

	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
	"github.com/nori-io/nori/pkg/nori/domain/errors"
)

type Registry struct {
	mx      sync.Mutex
	plugins []*entity.Plugin
}

func (r *Registry) Add(p *entity.Plugin) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	if p == nil {
		// todo: custom error
		return fmt.Errorf("plugin cannot be nil")
	}

	item := r.getByID(p.Meta().GetID())
	if item != nil {
		return errors.AlreadyExists{
			ID: p.Meta().GetID(),
		}
	}
	r.plugins = append(r.plugins, p)
	return nil
}

func (r *Registry) Remove(id meta.ID) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for i, item := range r.plugins {
		if item.Meta().GetID() == id {
			r.plugins = append(r.plugins[:i], r.plugins[:i+1]...)
		}
	}

	return nil
}

func (r *Registry) GetAll() []*entity.Plugin {
	return r.plugins
}

func (r *Registry) GetByID(id meta.ID) *entity.Plugin {
	r.mx.Lock()
	defer r.mx.Unlock()

	return r.getByID(id)
}

func (r *Registry) GetByInterface(i meta.Interface) []*entity.Plugin {
	r.mx.Lock()
	defer r.mx.Unlock()

	var plugins []*entity.Plugin
	for _, p := range r.plugins {
		if p.Meta().GetInterface().Equal(i) {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

// id: for future use
func (r *Registry) ResolveDependency(id meta.ID, d meta.Dependency) []*entity.Plugin {
	r.mx.Lock()
	defer r.mx.Unlock()

	var plugins []*entity.Plugin

	if d.Name() == "" {
		return plugins
	}
	if d.Constraint() == "" {
		return plugins
	}

	for _, p := range r.plugins {
		if p.Meta().GetInterface().Name() != d.Name() {
			continue
		}
		// todo: compare interface version and dependency constraint

		plugins = append(plugins, p)
	}
	return plugins
}

// nori-io/common Registry interface
func (r *Registry) ID(id meta.ID) (interface{}, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for i := range r.plugins {
		if r.plugins[i].Meta().GetID().String() == id.String() {
			return r.plugins[i].Instance(), nil
		}
	}
	return nil, errors.NotFound{ID: id}
}

// nori-io/common Registry interface
func (r *Registry) Interface(i meta.Interface) (interface{}, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, p := range r.plugins {
		if p.Meta().GetInterface().Name() != i.Name() {
			continue
		}
		// todo: compare interface versions

		return p.Instance(), nil
	}
	return nil, errors.InterfaceNotFound{Interface: i}
}

// nori-io/common Registry interface
func (r *Registry) Resolve(dep meta.Dependency) (interface{}, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, p := range r.plugins {
		if p.Meta().GetInterface().Name() != dep.Name() {
			continue
		}
		// todo: compare interface version and dependency constraint

		return p.Instance(), nil
	}
	return nil, errors.DependencyNotFound{Dependency: dep}
}

func (r *Registry) getByID(id meta.ID) *entity.Plugin {
	for _, p := range r.plugins {
		if p.Meta().GetID().String() == id.String() {
			return p
		}
	}
	return nil
}
