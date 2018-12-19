package plugins

import (
	"fmt"

	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/dependency"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/secure2work/nori/core/plugins/plugin"
	"github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph/simple"
)

type RegistryManager interface {
	Registry() plugin.Registry

	Add(p plugin.Plugin) error
	GetInterface(alias meta.Interface) interface{}
	Remove(p plugin.Plugin)
	Resolve(dep meta.Dependency) plugin.Plugin

	OrderedPluginList() (PluginList, error)
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

func (r registryManager) Add(p plugin.Plugin) error {
	// add plugin
	id := p.Meta().Id()
	r.plugins.Add(p)
	// add alias (non-Custom interfaces only)
	if p.Meta().GetInterface() != meta.Custom {
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
	}

	return nil
}

func (r registryManager) Remove(p plugin.Plugin) {
	r.plugins.Remove(p)

	// remove alias for non Custom interface
	if p.Meta().GetInterface() != meta.Custom {
		// @todo (?) replace with another plugin from the plugins list
		delete(r.interfaces, p.Meta().GetInterface())
	}
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

func (r registryManager) GetInterface(alias meta.Interface) interface{} {
	id, ok := r.interfaces[alias]
	if !ok {
		return nil
	}

	plugin, err := r.plugins.Find(id)
	if err != nil {
		return nil
	}

	return plugin.Instance()
}

func (r registryManager) Registry() plugin.Registry {
	return r.registry
}

func (r registryManager) OrderedPluginList() (PluginList, error) {
	graph := simple.NewDirectedGraph()

	nodes := map[meta.ID]dependency.Node{}

	for idx, p := range *r.plugins {
		node := dependency.NewNode(int64(idx), p.Meta().Id())
		graph.AddNode(node)
		nodes[p.Meta().Id()] = node
	}

	for _, p := range *r.plugins {
		n1 := nodes[p.Meta().Id()]
		for _, d := range p.Meta().GetDependencies() {
			var plug plugin.Plugin
			if d.Interface != meta.Custom {
				ifaceID, ok := r.interfaces[d.Interface]
				if !ok {
					return PluginList{}, fmt.Errorf("can't find dependency %s", d.ID)
				}
				plug, _ = r.plugins.Find(ifaceID)
				if plug == nil {
					return PluginList{}, fmt.Errorf("can't find dependency %s", d.ID)
				}
			} else {
				plug = r.plugins.Resolve(d)
			}
			if plug == nil {
				return PluginList{}, fmt.Errorf("can't find dependency %s", d.ID)
			}
			n2, ok := nodes[plug.Meta().Id()]
			if !ok {
				return PluginList{}, fmt.Errorf("can't find dependency %s", d.ID)
			}
			graph.SetEdge(simple.Edge{
				F: n1,
				T: n2,
			})
		}
	}

	ordered, err := dependency.Sort(graph)
	if err != nil {
		return PluginList{}, err
	}

	var list PluginList
	for _, n := range ordered {
		item, err := r.plugins.Find((*n).(dependency.PNode).PID())
		if err != nil {
			return PluginList{}, err
		}
		list.Add(item)
	}

	return list, nil
}
