package registry

import (
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/pkg/errors"
)

type Service struct {
	PluginRepository repository.PluginRepository
}

// implements github.com/nori-io/nori-common/v2/plugin/Plugin
func (r *Service) ID(id meta.ID) (interface{}, error) {
	p := r.PluginRepository.FindByID(id)
	if p == nil {
		return nil, errors.NotFound{ID: id}
	}
	return p.Instance(), nil
}

// implements github.com/nori-io/nori-common/v2/plugin/Plugin
func (r *Service) Interface(i meta.Interface) (interface{}, error) {
	for _, p := range r.PluginRepository.FindAll() {
		if p.Meta().GetInterface().Equal(i) {
			return p.Instance(), nil
		}
	}
	return nil, errors.InterfaceNotFound{
		Interface: i,
	}
}

// implements github.com/nori-io/nori-common/v2/plugin/Plugin
func (r *Service) Resolve(dep meta.Dependency) (interface{}, error) {
	cons, err := dep.GetConstraint()
	if err != nil {
		return nil, err
	}

	for _, p := range r.PluginRepository.FindAll() {
		if dep.Interface.Name() != p.Meta().GetInterface().Name() {
			continue
		}

		v, _ := p.Meta().Id().GetVersion()

		if cons.Check(v) {
			return p.Instance(), nil
		}
	}
	return nil, errors.DependencyNotFound{
		Dependency: dep,
	}
}
