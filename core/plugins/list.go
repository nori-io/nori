package plugins

import (
	"github.com/secure2work/nori/core/plugins/errors"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/secure2work/nori/core/plugins/plugin"
)

type PluginList []plugin.Plugin

func (pl PluginList) Find(id meta.ID) (plugin.Plugin, error) {
	for _, p := range pl {
		if p.Meta().Id() == id {
			return p, nil
		}
	}
	return nil, errors.NotFound{
		ID: id,
	}
}

func (pl PluginList) FindByPluginID(id meta.PluginID) PluginList {
	list := PluginList{}
	for _, p := range pl {
		if p.Meta().Id().ID == id {
			list = append(list, p)
		}
	}
	return list
}

func (pl PluginList) Resolve(dependency meta.Dependency) plugin.Plugin {
	cons, err := dependency.GetConstraint()
	if err != nil {
		return nil
	}

	for _, p := range pl {
		if dependency.ID != p.Meta().Id().ID {
			continue
		}

		v, _ := p.Meta().Id().GetVersion()

		if cons.Check(v) {
			return p
		}
	}
	return nil
}

func (pl *PluginList) Add(p plugin.Plugin) {
	if p, _ := pl.Find(p.Meta().Id()); p != nil {
		return
	}
	*pl = append(*pl, p)
}

func (pl *PluginList) Remove(p plugin.Plugin) {
	for i, v := range *pl {
		if v == p {
			*pl = append((*pl)[:i], (*pl)[i+1:]...)
		}
	}
}
