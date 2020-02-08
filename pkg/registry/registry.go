/*
Copyright 2019-2020 The Nori Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package registry

import (
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori-common/v2/plugin"
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

// implements nori-common.Registry.Resolve()
func (r *registry) Resolve(dep meta.Dependency) (interface{}, error) {
	p, err := r.plugins.Resolve(dep)
	if err != nil {
		return nil, err
	}
	return p.Instance(), nil
}
