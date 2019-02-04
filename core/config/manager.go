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

package config

import (
	commonCfg "github.com/nori-io/nori-common/config"
	"github.com/nori-io/nori-common/interfaces"
	"github.com/nori-io/nori-common/meta"
)

type manager struct {
	configs map[meta.ID]*[]commonCfg.Variable
	config  interfaces.Config
}

func NewManager(config interfaces.Config) commonCfg.Manager {
	m := new(manager)
	m.configs = make(map[meta.ID]*[]commonCfg.Variable)
	m.config = config
	return m
}

func (m *manager) Register(meta meta.Meta) commonCfg.Config {
	cfgs := make([]commonCfg.Variable, 0)
	m.configs[meta.Id()] = &cfgs
	return &config{
		cfgs:   &cfgs,
		config: m.config,
	}
}

func (m *manager) PluginVariables(id meta.ID) []commonCfg.Variable {
	vars, ok := m.configs[id]
	if !ok {
		return []commonCfg.Variable{}
	}
	return *vars
}

type config struct {
	cfgs   *[]commonCfg.Variable
	config interfaces.Config
}

func (c *config) Bool(key string, desc string) commonCfg.Bool {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
	})
	return func() bool {
		return c.config.Bool(key)
	}
}

func (c *config) Float(key string, desc string) commonCfg.Float {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
	})
	return func() float64 {
		return c.config.Float(key)
	}
}

func (c *config) Int(key string, desc string) commonCfg.Int {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
	})
	return func() int {
		return c.config.Int(key)
	}
}

func (c *config) UInt(key string, desc string) commonCfg.UInt {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
	})
	return func() uint {
		return c.config.UInt(key)
	}
}

func (c *config) Slice(key, delimiter string, desc string) commonCfg.Slice {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
	})
	return func() []interface{} {
		return c.config.Slice(key, delimiter)
	}
}

func (c *config) String(key string, desc string) commonCfg.String {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
	})
	return func() string {
		return c.config.String(key)
	}
}

func (c *config) StringMap(key string, desc string) commonCfg.StringMap {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
	})
	return func() map[string]interface{} {
		return c.config.StringMap(key)
	}
}
