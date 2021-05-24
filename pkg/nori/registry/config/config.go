package config

import (
	go_config "github.com/cheebo/go-config"
	"github.com/nori-io/common/v5/pkg/domain/config"
	enum "github.com/nori-io/common/v5/pkg/domain/enum/config"
	"github.com/nori-io/common/v5/pkg/domain/registry"
)

type Config struct {
	config    go_config.Fields
	variables *[]registry.Variable
}

func (c *Config) Bool(key string, desc string) config.Bool {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Bool,
	})
	return func() bool {
		return c.config.Bool(key)
	}
}

func (c *Config) Float(key string, desc string) config.Float {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Float,
	})
	return func() float64 {
		return c.config.Float(key)
	}
}

func (c *Config) Int(key string, desc string) config.Int {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Int,
	})
	return func() int {
		return c.config.Int(key)
	}
}

func (c *Config) Int8(key string, desc string) config.Int8 {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Int8,
	})
	return func() int8 {
		return c.config.Int8(key)
	}
}

func (c *Config) Int32(key string, desc string) config.Int32 {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Int32,
	})
	return func() int32 {
		return c.config.Int32(key)
	}
}

func (c *Config) Int64(key string, desc string) config.Int64 {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Int64,
	})
	return func() int64 {
		return c.config.Int64(key)
	}
}

func (c *Config) UInt(key string, desc string) config.UInt {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.UInt,
	})
	return func() uint {
		return c.config.UInt(key)
	}
}

func (c *Config) UInt32(key string, desc string) config.UInt32 {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.UInt32,
	})
	return func() uint32 {
		return c.config.UInt32(key)
	}
}

func (c *Config) UInt64(key string, desc string) config.UInt64 {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.UInt64,
	})
	return func() uint64 {
		return c.config.UInt64(key)
	}
}

func (c *Config) Slice(key, desc string) config.Slice {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Slice,
	})
	return func() []interface{} {
		return c.config.Slice(key)
	}
}

func (c *Config) SliceInt(key, desc string) config.SliceInt {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Slice,
	})
	return func() []int {
		return c.config.SliceInt(key)
	}
}

func (c *Config) SliceString(key, desc string) config.SliceString {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Slice,
	})
	return func() []string {
		return c.config.SliceString(key)
	}
}

func (c *Config) String(key string, desc string) config.String {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.String,
	})
	return func() string {
		return c.config.String(key)
	}
}

func (c *Config) StringMap(key string, desc string) config.StringMap {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Map,
	})
	return func() map[string]interface{} {
		return c.config.StringMap(key)
	}
}

func (c *Config) StringMapInt(key string, desc string) config.StringMapInt {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Map,
	})
	return func() map[string]int {
		return c.config.StringMapInt(key)
	}
}

func (c *Config) StringMapSliceString(key string, desc string) config.StringMapSliceString {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Map,
	})
	return func() map[string][]string {
		return c.config.StringMapSliceString(key)
	}
}

func (c *Config) StringMapString(key string, desc string) config.StringMapString {
	*(c.variables) = append(*(c.variables), registry.Variable{
		Name:        key,
		Description: desc,
		Type:        enum.Map,
	})
	return func() map[string]string {
		return c.config.StringMapString(key)
	}
}

func (c *Config) SetDefault(key string, val interface{}) {
	//c.config.SetDefault(key, val)
}
