package config

import (
	"github.com/secure2work/nori/core/entities"
	"github.com/secure2work/nori/core/plugins/interfaces"
)

type manager struct {
	configs map[string][]ConfigVariable
	config  interfaces.Config
}

func NewConfigManager(config interfaces.Config) interfaces.ConfigManager {
	m := new(manager)
	m.configs = make(map[string][]ConfigVariable)
	m.config = config
	return m
}

func (m *manager) PluginRegister(meta entities.PluginMeta) interfaces.Manager {
	cfgs := make([]ConfigVariable, 0)
	m.configs[meta.GetId()] = cfgs
	return &configs{
		cfgs:   cfgs,
		config: m.config,
	}
}

type ConfigVariable struct {
	Name        string
	Description string
}

type configs struct {
	cfgs   []ConfigVariable
	config interfaces.Config
}

func (c *configs) Bool(key string, desc string) interfaces.Bool {
	c.cfgs = append(c.cfgs, ConfigVariable{
		Name:        key,
		Description: desc,
	})
	return func() bool {
		return c.config.Bool(key)
	}
}

func (c *configs) Float(key string, desc string) interfaces.Float {
	c.cfgs = append(c.cfgs, ConfigVariable{
		Name:        key,
		Description: desc,
	})
	return func() float64 {
		return c.config.Float(key)
	}
}

func (c *configs) Int(key string, desc string) interfaces.Int {
	c.cfgs = append(c.cfgs, ConfigVariable{
		Name:        key,
		Description: desc,
	})
	return func() int {
		return c.config.Int(key)
	}
}

func (c *configs) UInt(key string, desc string) interfaces.UInt {
	c.cfgs = append(c.cfgs, ConfigVariable{
		Name:        key,
		Description: desc,
	})
	return func() uint {
		return c.config.UInt(key)
	}
}

func (c *configs) Slice(key, delimiter string, desc string) interfaces.Slice {
	c.cfgs = append(c.cfgs, ConfigVariable{
		Name:        key,
		Description: desc,
	})
	return func() []interface{} {
		return c.config.Slice(key, delimiter)
	}
}

func (c *configs) String(key string, desc string) interfaces.String {
	c.cfgs = append(c.cfgs, ConfigVariable{
		Name:        key,
		Description: desc,
	})
	return func() string {
		return c.config.String(key)
	}
}

func (c *configs) StringMap(key string, desc string) interfaces.StringMap {
	c.cfgs = append(c.cfgs, ConfigVariable{
		Name:        key,
		Description: desc,
	})
	return func() map[string]interface{} {
		return c.config.StringMap(key)
	}
}
