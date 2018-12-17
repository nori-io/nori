package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/secure2work/nori/plugin"
	"github.com/sirupsen/logrus"

	"context"
)

func TestNewRegistry(t *testing.T) {
	a := assert.New(t)

	r := plugin.NewRegistry(nil, logrus.New())

	id := plugin.ID{
		ID:      "nori/test",
		Version: "1.1",
	}

	p := plugguble{}
	p.meta = plugin.MetaData{
		ID:           id,
		Dependencies: []plugin.Dependency{},
		Interface:    plugin.Custom,
		Core: plugin.Core{
			VersionConstraint: "~1.0",
		},
	}

	r.Add(p)

	dep := plugin.Dependency{
		ID:         "nori/test",
		Constraint: ">=1.0",
	}

	found := r.Resolve(dep)

	a.NotNil(found)
}

type plugguble struct {
	meta plugin.MetaData
}

func (p plugguble) GetInstance() interface{} {
	return p
}
func (p plugguble) GetMeta() plugin.Meta {
	return p.meta
}

func (p plugguble) Init(ctx context.Context, config plugin.ConfigManager) error {
	return nil
}

func (p plugguble) Install(ctx context.Context, registry plugin.Registry) error {
	return nil
}

func (p plugguble) Start(ctx context.Context, registry plugin.Registry) error {
	return nil
}

func (p plugguble) Stop(ctx context.Context, registry plugin.Registry) error {
	return nil
}

func (p plugguble) UnInstall(ctx context.Context, registry plugin.Registry) error {
	return nil
}
