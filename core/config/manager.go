package config

import (
	"github.com/secure2work/nori/core/plugins/interfaces"
	"github.com/secure2work/nori/core/plugins/meta"
)

type Manager interface {
	Register(meta.Meta) Config
}

type Config interface {
	Bool(key string, desc string) Bool
	Float(key string, desc string) Float
	Int(key string, desc string) Int
	UInt(key string, desc string) UInt
	Slice(key, delimiter string, desc string) Slice
	String(key string, desc string) String
	StringMap(key string, desc string) StringMap
}

type (
	Bool      func() bool
	Float     func() float64
	Int       func() int
	UInt      func() uint
	Slice     func() []interface{}
	String    func() string
	StringMap func() map[string]interface{}
)

type Variable struct {
	Name        string
	Description string
}

type manager struct {
	configs map[meta.ID][]Variable
	config  interfaces.Config
}

func NewManager(config interfaces.Config) Manager {
	m := new(manager)
	m.configs = make(map[meta.ID][]Variable)
	m.config = config
	return m
}

func (m *manager) Register(meta meta.Meta) Config {
	cfgs := make([]Variable, 0)
	m.configs[meta.Id()] = cfgs
	return &config{
		cfgs:   cfgs,
		config: m.config,
	}
}

type config struct {
	cfgs   []Variable
	config interfaces.Config
}

func (c *config) Bool(key string, desc string) Bool {
	c.cfgs = append(c.cfgs, Variable{
		Name:        key,
		Description: desc,
	})
	return func() bool {
		return c.config.Bool(key)
	}
}

func (c *config) Float(key string, desc string) Float {
	c.cfgs = append(c.cfgs, Variable{
		Name:        key,
		Description: desc,
	})
	return func() float64 {
		return c.config.Float(key)
	}
}

func (c *config) Int(key string, desc string) Int {
	c.cfgs = append(c.cfgs, Variable{
		Name:        key,
		Description: desc,
	})
	return func() int {
		return c.config.Int(key)
	}
}

func (c *config) UInt(key string, desc string) UInt {
	c.cfgs = append(c.cfgs, Variable{
		Name:        key,
		Description: desc,
	})
	return func() uint {
		return c.config.UInt(key)
	}
}

func (c *config) Slice(key, delimiter string, desc string) Slice {
	c.cfgs = append(c.cfgs, Variable{
		Name:        key,
		Description: desc,
	})
	return func() []interface{} {
		return c.config.Slice(key, delimiter)
	}
}

func (c *config) String(key string, desc string) String {
	c.cfgs = append(c.cfgs, Variable{
		Name:        key,
		Description: desc,
	})
	return func() string {
		return c.config.String(key)
	}
}

func (c *config) StringMap(key string, desc string) StringMap {
	c.cfgs = append(c.cfgs, Variable{
		Name:        key,
		Description: desc,
	})
	return func() map[string]interface{} {
		return c.config.StringMap(key)
	}
}
