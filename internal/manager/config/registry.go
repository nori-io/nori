package config

import (
	go_config "github.com/cheebo/go-config"
	commonCfg "github.com/nori-io/nori-common/v2/config"
	"github.com/nori-io/nori-common/v2/config/types"
)

type registry struct {
	cfgs   *[]commonCfg.Variable
	config go_config.Fields
}

func (c *registry) Bool(key string, desc string) commonCfg.Bool {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Bool,
	})
	return func() bool {
		return c.config.Bool(key)
	}
}

func (c *registry) Float(key string, desc string) commonCfg.Float {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Float,
	})
	return func() float64 {
		return c.config.Float(key)
	}
}

func (c *registry) Int(key string, desc string) commonCfg.Int {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Int,
	})
	return func() int {
		return c.config.Int(key)
	}
}

func (c *registry) Int8(key string, desc string) commonCfg.Int8 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Int8,
	})
	return func() int8 {
		return c.config.Int8(key)
	}
}

func (c *registry) Int32(key string, desc string) commonCfg.Int32 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Int32,
	})
	return func() int32 {
		return c.config.Int32(key)
	}
}

func (c *registry) Int64(key string, desc string) commonCfg.Int64 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Int64,
	})
	return func() int64 {
		return c.config.Int64(key)
	}
}

func (c *registry) UInt(key string, desc string) commonCfg.UInt {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.UInt,
	})
	return func() uint {
		return c.config.UInt(key)
	}
}

func (c *registry) UInt32(key string, desc string) commonCfg.UInt32 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.UInt32,
	})
	return func() uint32 {
		return c.config.UInt32(key)
	}
}

func (c *registry) UInt64(key string, desc string) commonCfg.UInt64 {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.UInt64,
	})
	return func() uint64 {
		return c.config.UInt64(key)
	}
}

func (c *registry) Slice(key, desc string) commonCfg.Slice {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Slice,
	})
	return func() []interface{} {
		return c.config.Slice(key)
	}
}

func (c *registry) SliceInt(key, desc string) commonCfg.SliceInt {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Slice,
	})
	return func() []int {
		return c.config.SliceInt(key)
	}
}

func (c *registry) SliceString(key, desc string) commonCfg.SliceString {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Slice,
	})
	return func() []string {
		return c.config.SliceString(key)
	}
}

func (c *registry) String(key string, desc string) commonCfg.String {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.String,
	})
	return func() string {
		return c.config.String(key)
	}
}

func (c *registry) StringMap(key string, desc string) commonCfg.StringMap {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Map,
	})
	return func() map[string]interface{} {
		return c.config.StringMap(key)
	}
}

func (c *registry) StringMapInt(key string, desc string) commonCfg.StringMapInt {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Map,
	})
	return func() map[string]int {
		return c.config.StringMapInt(key)
	}
}

func (c *registry) StringMapSliceString(key string, desc string) commonCfg.StringMapSliceString {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Map,
	})
	return func() map[string][]string {
		return c.config.StringMapSliceString(key)
	}
}

func (c *registry) StringMapString(key string, desc string) commonCfg.StringMapString {
	*(c.cfgs) = append(*(c.cfgs), commonCfg.Variable{
		Name:        key,
		Description: desc,
		Type:        types.Map,
	})
	return func() map[string]string {
		return c.config.StringMapString(key)
	}
}
