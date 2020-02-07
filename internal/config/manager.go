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
	"github.com/cheebo/go-config"
	commonCfg "github.com/nori-io/nori-common/v2/config"
	"github.com/nori-io/nori-common/v2/config/types"
	"github.com/nori-io/nori-common/v2/meta"
)

type manager struct {
	configs map[meta.ID]*[]commonCfg.Variable
	config  go_config.Fields
}

func NewManager(config go_config.Fields) commonCfg.Manager {
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
	config go_config.Fields
}

func (c *config) Bool(key string, desc string) commonCfg.Bool {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Bool,
	})
	return func() bool {
		return c.config.Bool(key)
	}
}

func (c *config) Float(key string, desc string) commonCfg.Float {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Float,
	})
	return func() float64 {
		return c.config.Float(key)
	}
}

func (c *config) Int(key string, desc string) commonCfg.Int {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Int,
	})
	return func() int {
		return c.config.Int(key)
	}
}

func (c *config) Int8(key string, desc string) commonCfg.Int8 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Int8,
	})
	return func() int8 {
		return c.config.Int8(key)
	}
}

func (c *config) Int32(key string, desc string) commonCfg.Int32 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Int32,
	})
	return func() int32 {
		return c.config.Int32(key)
	}
}

func (c *config) Int64(key string, desc string) commonCfg.Int64 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Int64,
	})
	return func() int64 {
		return c.config.Int64(key)
	}
}

func (c *config) UInt(key string, desc string) commonCfg.UInt {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.UInt,
	})
	return func() uint {
		return c.config.UInt(key)
	}
}

func (c *config) UInt32(key string, desc string) commonCfg.UInt32 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.UInt32,
	})
	return func() uint32 {
		return c.config.UInt32(key)
	}
}

func (c *config) UInt64(key string, desc string) commonCfg.UInt64 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.UInt64,
	})
	return func() uint64 {
		return c.config.UInt64(key)
	}
}

func (c *config) Slice(key, desc string) commonCfg.Slice {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Slice,
	})
	return func() []interface{} {
		return c.config.Slice(key)
	}
}

func (c *config) SliceInt(key, desc string) commonCfg.SliceInt {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Slice,
	})
	return func() []int {
		return c.config.SliceInt(key)
	}
}

func (c *config) SliceString(key, desc string) commonCfg.SliceString {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Slice,
	})
	return func() []string {
		return c.config.SliceString(key)
	}
}

func (c *config) String(key string, desc string) commonCfg.String {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.String,
	})
	return func() string {
		return c.config.String(key)
	}
}

func (c *config) StringMap(key string, desc string) commonCfg.StringMap {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Map,
	})
	return func() map[string]interface{} {
		return c.config.StringMap(key)
	}
}

func (c *config) StringMapInt(key string, desc string) commonCfg.StringMapInt {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Map,
	})
	return func() map[string]int {
		return c.config.StringMapInt(key)
	}
}

func (c *config) StringMapSliceString(key string, desc string) commonCfg.StringMapSliceString {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Map,
	})
	return func() map[string][]string {
		return c.config.StringMapSliceString(key)
	}
}

func (c *config) StringMapString(key string, desc string) commonCfg.StringMapString {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Map,
	})
	return func() map[string]string {
		return c.config.StringMapString(key)
	}
}
