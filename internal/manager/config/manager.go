/*
Copyright 2018-2020 The Nori Authors.
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

package config

import (
	"github.com/cheebo/go-config"
	"github.com/nori-io/nori-common/v2/config"
	"github.com/nori-io/nori-common/v2/meta"
)

type manager struct {
	configs map[meta.ID]*[]config.Variable
	config  go_config.Fields
}

func (m *manager) Register(id meta.ID) config.Config {
	vars := make([]config.Variable, 0)
	m.configs[id] = &vars
	return &registry{
		cfgs:   &vars,
		config: m.config,
	}
}

func (m *manager) PluginVariables(id meta.ID) []config.Variable {
	vars, ok := m.configs[id]
	if !ok {
		return []config.Variable{}
	}
	return *vars
}
