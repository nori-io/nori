package interfaces

import (
	"github.com/secure2work/nori/core/entities"
)

type ConfigManager interface {
	PluginRegister(entities.PluginMeta) Manager
}

type Manager interface {
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
