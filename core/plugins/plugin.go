package plugins

import (
	"context"

	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/meta"
)

type Plugin interface {
	Meta() meta.Meta
	Instance() interface{}
	Init(ctx context.Context, config config.Manager) error
	Start(ctx context.Context, registry Registry) error
	Stop(ctx context.Context, registry Registry) error
}

type Installable interface {
	Install(ctx context.Context, registry Registry) error
	UnInstall(ctx context.Context, registry Registry) error
}

type PluginList []Plugin

func (pl PluginList) Find(id meta.ID) (Plugin, error) {
	for _, p := range pl {
		if p.Meta().Id() == id {
			return p, nil
		}
	}
	return nil, NotFound{
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

func (pl PluginList) Resolve(dependency meta.Dependency) Plugin {
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

func (pl *PluginList) Add(p Plugin) {
	if p, _ := pl.Find(p.Meta().Id()); p != nil {
		return
	}
	*pl = append(*pl, p)
}

func (pl *PluginList) Remove(p Plugin) {
	for i, v := range *pl {
		if v == p {
			*pl = append((*pl)[:i], (*pl)[i+1:]...)
		}
	}
}
