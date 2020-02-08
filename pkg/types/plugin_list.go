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

package types

import (
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori-common/v2/plugin"
	"github.com/nori-io/nori/pkg/errors"
)

type PluginList []plugin.Plugin

func (pl *PluginList) Add(p plugin.Plugin) error {
	if p, _ := pl.ID(p.Meta().Id()); p != nil {
		return errors.AlreadyExists{
			ID: p.Meta().Id(),
		}
	}
	*pl = append(*pl, p)
	return nil
}

func (pl *PluginList) ID(id meta.ID) (plugin.Plugin, error) {
	for _, p := range *pl {
		if p.Meta().Id() == id {
			return p, nil
		}
	}
	return nil, errors.NotFound{
		ID: id,
	}
}

func (pl *PluginList) Interface(i meta.Interface) (plugin.Plugin, error) {
	for _, p := range *pl {
		if p.Meta().GetInterface().Equal(i) {
			return p, nil
		}
	}
	return nil, errors.InterfaceNotFound{Interface: i}
}

func (pl *PluginList) Remove(p plugin.Plugin) {
	for i, v := range *pl {
		if v == p {
			*pl = append((*pl)[:i], (*pl)[i+1:]...)
		}
	}
}

func (pl *PluginList) Resolve(dep meta.Dependency) (plugin.Plugin, error) {
	cons, err := dep.GetConstraint()
	if err != nil {
		return nil, err
	}

	for _, p := range *pl {
		if dep.Interface != p.Meta().GetInterface() {
			continue
		}

		v, _ := p.Meta().Id().GetVersion()

		if cons.Check(v) {
			return p, nil
		}
	}
	return nil, errors.DependencyNotFound{Dependency: dep}
}
