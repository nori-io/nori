package plugins

import (
	"github.com/secure2work/nori/core/plugins/errors"

	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/secure2work/nori/core/plugins/plugin"
	"github.com/sirupsen/logrus"
)

type RegistryManager interface {
	//Registry() plugin.Registry

	Add(p plugin.Plugin) error
	Get(id meta.ID) (plugin.Plugin, error)
	GetInterface(alias meta.Interface) (interface{}, error)
	Plugins() PluginList
	Resolve(dep meta.Dependency) plugin.Plugin
	Remove(p plugin.Plugin)
}

type registryManager struct {
	log           *logrus.Logger
	plugins       *PluginList
	interfaces    map[meta.Interface]meta.ID
	configManager config.Manager
	registry      plugin.Registry
}

func NewRegistryManager(cm config.Manager, logger *logrus.Logger) RegistryManager {
	rm := &registryManager{
		log:        logger,
		plugins:    &PluginList{},
		interfaces: map[meta.Interface]meta.ID{},
	}
	rm.registry = NewRegistry(rm, cm, logger)
	return rm
}

//func (r registryManager) Registry() plugin.Registry {
//	return r.registry
//}

func (r registryManager) Add(p plugin.Plugin) error {
	// add plugin
	id := p.Meta().Id()
	r.plugins.Add(p)

	if p.Meta().GetInterface() == meta.Custom {
		return nil
	}

	// add alias (non-Custom interfaces only)
	// 1. if alias is empty - fill it with plugin, otherwise
	// 2. if alias.ID equal to plugin.ID,
	// then take plugin that has greater version, otherwise
	// 3. fill alias with new plugin
	alias, ok := r.interfaces[p.Meta().GetInterface()]
	if !ok {
		r.interfaces[p.Meta().GetInterface()] = id
	} else {
		if id.ID != alias.ID {
			r.interfaces[p.Meta().GetInterface()] = id
		} else {
			newVer, err := id.GetVersion()
			if err != nil {
				r.log.Error(err)
				return err
			}
			curVer, err := alias.GetVersion()
			if err != nil {
				r.log.Error(err)
				return err
			}
			if newVer.GreaterThan(curVer) {
				r.interfaces[p.Meta().GetInterface()] = id
			}
		}
	}

	return nil
}

func (r registryManager) Get(id meta.ID) (plugin.Plugin, error) {
	return r.plugins.Find(id)
}

func (r registryManager) GetInterface(alias meta.Interface) (interface{}, error) {
	id, ok := r.interfaces[alias]
	if !ok {
		return nil, errors.InterfaceNotFound{
			Interface: alias,
		}
	}

	p, err := r.plugins.Find(id)
	if err != nil {
		return nil, errors.NotFound{
			ID: id,
		}
	}

	return p.Instance(), nil
}

func (r registryManager) Remove(p plugin.Plugin) {
	r.plugins.Remove(p)

	// remove alias for non Custom interface
	if p.Meta().GetInterface() != meta.Custom {
		// @todo (?) replace with another plugin from the plugins list
		delete(r.interfaces, p.Meta().GetInterface())
	}
}

func (r registryManager) Plugins() PluginList {
	return *r.plugins
}

func (r registryManager) Resolve(dep meta.Dependency) plugin.Plugin {
	if dep.Interface != meta.Custom {
		id, ok := r.interfaces[dep.Interface]
		if !ok {
			return nil
		}
		plug, err := r.plugins.Find(id)
		if err != nil {
			return nil
		}
		return plug
	}
	return r.plugins.Resolve(dep)
}
