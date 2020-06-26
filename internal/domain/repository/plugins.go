package repository

import (
	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori/internal/domain/entity"
)

type PluginRepository interface {
	// Register adds plugin to plugin set
	Register(p *entity.Plugin) error
	// UnRegister removes plugin from plugin set
	UnRegister(p *entity.Plugin) error

	// FindAll returns all registered plugins
	FindAll() []*entity.Plugin
	// FindByID returns plugin by specific id
	FindByID(id meta.ID) *entity.Plugin
	// FindByIDs returns set of plugins matched by ids
	FindByIDs(ids []meta.ID) []*entity.Plugin
	// FindByInterface returns all plugin that implement interface i
	FindByInterface(i meta.Interface) []*entity.Plugin

	// Resolve returns plugin that fits given dependency spec
	Resolve(dependency meta.Dependency) (*entity.Plugin, error)
	// FindDependent returns dependent plugins
	FindDependent(id meta.ID) []*entity.Plugin

	// FindInstallable returns all plugins that implement Installable interface
	FindInstallable() []*entity.Plugin
	// FindHooks returns all plugins that implement Hook interface
	FindHooks() []*entity.Plugin

	// Tree returns plugin execution tree as an array
	Tree() ([]*entity.Plugin, error)
}
