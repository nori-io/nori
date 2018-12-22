package plugins_test

import (
	"testing"

	noriPlugin "github.com/secure2work/nori/core/plugins/plugin"

	"github.com/secure2work/nori/core/config"

	"github.com/stretchr/testify/assert"

	"github.com/secure2work/nori/core/plugins"
	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/sirupsen/logrus"

	"context"
)

func TestNewRegistry(t *testing.T) {
	a := assert.New(t)

	r := plugins.NewRegistryManager(nil, logrus.New())

	id := meta.ID{
		ID:      "nori/test",
		Version: "1.1",
	}

	p := plugguble{}
	p.meta = meta.Data{
		ID:           id,
		Dependencies: []meta.Dependency{},
		Interface:    meta.Custom,
		Core: meta.Core{
			VersionConstraint: "~1.0",
		},
	}

	r.Add(p)

	dep := meta.Dependency{
		ID:         "nori/test",
		Constraint: ">=1.0",
	}

	found, err := r.Registry().Resolve(dep)

	a.NoError(err)

	a.NotNil(found)
}

func TestRegistryManager_OrderedPluginList_PluginsOnly(t *testing.T) {
	a := assert.New(t)

	r := plugins.NewRegistryManager(nil, logrus.New())

	one := plugguble{}
	one.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/core",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{},
		Interface:    meta.Custom,
	}
	err := r.Add(one)
	a.NoError(err)

	two := plugguble{}
	two.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/api",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{
			meta.Dependency{
				ID:         meta.PluginID("nori/core"),
				Constraint: "1.1",
			},
		},
		Interface: meta.Custom,
	}
	r.Add(two)

	list, err := r.OrderedPluginList()
	a.NoError(err)
	a.Len(list, 2)

	if len(list) == 2 {
		r1 := list[0]
		r2 := list[1]
		a.Equal("nori/core", string(r1.Meta().Id().ID))
		a.Equal("nori/api", string(r2.Meta().Id().ID))
	}
}

func TestRegistryManager_OrderedPluginList_Interfaces(t *testing.T) {
	a := assert.New(t)

	r := plugins.NewRegistryManager(nil, logrus.New())

	two := plugguble{}
	two.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/cache",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{},
		Interface:    meta.Cache,
	}
	r.Add(two)

	three := plugguble{}
	three.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/api",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{
			meta.Cache.Dependency(),
		},
		Interface: meta.Custom,
	}
	r.Add(three)

	list, err := r.OrderedPluginList()
	a.NoError(err)
	a.Len(list, 2)

	if len(list) == 2 {
		r1 := list[0]
		r2 := list[1]
		a.Equal("nori/cache", string(r1.Meta().Id().ID))
		a.Equal("nori/api", string(r2.Meta().Id().ID))
	}
}

func TestRegistryManager_OrderedPluginList_InterfaceAndPlugin(t *testing.T) {
	a := assert.New(t)

	r := plugins.NewRegistryManager(nil, logrus.New())

	one := plugguble{}
	one.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/dummy",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{
			meta.Dependency{
				ID:         "nori/api",
				Constraint: ">=1.0",
			},
		},
		Interface: meta.Custom,
	}
	r.Add(one)

	two := plugguble{}
	two.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/cache",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{},
		Interface:    meta.Cache,
	}
	r.Add(two)

	three := plugguble{}
	three.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/api",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{
			meta.Cache.Dependency(),
		},
		Interface: meta.Custom,
	}
	r.Add(three)

	list, err := r.OrderedPluginList()
	a.NoError(err)
	a.Len(list, 3)

	if len(list) == 3 {
		r1 := list[0]
		r2 := list[1]
		r3 := list[2]
		a.Equal("nori/cache", string(r1.Meta().Id().ID))
		a.Equal("nori/api", string(r2.Meta().Id().ID))
		a.Equal("nori/dummy", string(r3.Meta().Id().ID))
	}
}

func TestRegistryManager_OrderedPluginList_Error(t *testing.T) {
	a := assert.New(t)

	r := plugins.NewRegistryManager(nil, logrus.New())

	one := plugguble{}
	one.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/dummy",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{
			meta.Dependency{
				ID:         "nori/api",
				Constraint: ">=1.0",
			},
		},
		Interface: meta.Custom,
	}
	r.Add(one)

	two := plugguble{}
	two.meta = meta.Data{
		ID: meta.ID{
			ID:      "nori/cache",
			Version: "1.1",
		},
		Dependencies: []meta.Dependency{},
		Interface:    meta.Cache,
	}
	r.Add(two)

	_, err := r.OrderedPluginList()
	a.Error(err)
}

type plugguble struct {
	meta meta.Data
}

func (p plugguble) Instance() interface{} {
	return p
}
func (p plugguble) Meta() meta.Meta {
	return p.meta
}

func (p plugguble) Init(ctx context.Context, config config.Manager) error {
	return nil
}

func (p plugguble) Start(ctx context.Context, registry noriPlugin.Registry) error {
	return nil
}

func (p plugguble) Stop(ctx context.Context, registry noriPlugin.Registry) error {
	return nil
}
