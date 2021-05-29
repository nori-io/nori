package registry

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/plugin"
)

type Registry interface {
	Add(p plugin.Plugin) error
	Remove(p plugin.Plugin) error

	GetAll() []plugin.Plugin
	GetByID(id meta.ID) plugin.Plugin
	GetByInterface(i meta.Interface) []plugin.Plugin

	ResolveDependency(id meta.ID, d meta.Dependency) []plugin.Plugin

	// nori-io/common Registry interface
	ID(id meta.ID) (interface{}, error)
	Interface(i meta.Interface) (interface{}, error)
	Resolve(dep meta.Dependency) (interface{}, error)
}
