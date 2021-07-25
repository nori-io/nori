package registry

import (
	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
)

type Registry interface {
	Add(p *entity.Plugin) error
	Remove(id meta.ID) error

	GetAll() []*entity.Plugin
	GetByID(id meta.ID) *entity.Plugin
	GetByInterface(i meta.Interface) []*entity.Plugin

	ResolveDependency(id meta.ID, d meta.Dependency) []*entity.Plugin

	// nori-io/common Registry interface
	ID(id meta.ID) (interface{}, error)
	Interface(i meta.Interface) (interface{}, error)
	Resolve(dep meta.Dependency) (interface{}, error)
}
