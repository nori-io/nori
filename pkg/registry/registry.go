package registry

import (
	"github.com/nori-io/nori-common/logger"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori-common/plugin"
	"github.com/nori-io/nori/pkg/types"
)

type Registry interface {
	Add(p plugin.Plugin) error
	Get(id meta.ID) (plugin.Plugin, error)
	Remove(p plugin.Plugin)
	Plugins() []plugin.Plugin

	plugin.Registry
}

type registry struct {
	log     logger.Logger
	plugins types.PluginList
}

func NewRegistry(log logger.Logger) Registry {
	return &registry{
		log: log.With(logger.Field{
			Key:   "component",
			Value: "registry",
		}),
	}
}

func (r *registry) Add(p plugin.Plugin) error {
	return r.plugins.Add(p)
}

func (r *registry) Get(id meta.ID) (plugin.Plugin, error) {
	return r.plugins.ID(id)
}

func (r *registry) Remove(p plugin.Plugin) {
	r.plugins.Remove(p)
}

func (r *registry) Plugins() []plugin.Plugin {
	return r.plugins
}

// implements nori-common.Registry.ID()
func (r *registry) ID(id meta.ID) (interface{}, error) {
	p, err := r.plugins.ID(id)
	if err != nil {
		return nil, err
	}
	return p.Instance(), nil
}

// implements nori-common.Registry.Interface()
func (r *registry) Interface(i meta.Interface) (interface{}, error) {
	p, err := r.plugins.Interface(i)
	if err != nil {
		return nil, err
	}
	return p.Instance(), nil
}

// implements nori-common.Registry.Logger()
func (r *registry) Logger() logger.Logger {
	return r.log
}

// implements nori-common.Registry.Resolve()
func (r *registry) Resolve(dep meta.Dependency) (interface{}, error) {
	p, err := r.plugins.Resolve(dep)
	if err != nil {
		return nil, err
	}
	return p.Instance(), nil
}
