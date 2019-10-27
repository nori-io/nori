package plugins_test

import (
	"testing"

	"github.com/nori-io/nori/internal/log"

	noriPlugin "github.com/nori-io/nori-common/plugin"

	"github.com/nori-io/nori-common/config"

	"github.com/stretchr/testify/assert"

	"context"

	"github.com/nori-io/nori-common/meta"
)

func TestNewRegistry(t *testing.T) {
	a := assert.New(t)

	logger := log.New()
	r := NewRegistryManager(nil, logger)

	id := meta.ID{
		ID:      "nori/test",
		Version: "1.1",
	}

	p := plugguble{}
	p.meta = meta.Data{
		ID:           id,
		Dependencies: []meta.Dependency{},
		Interface:    meta.Interface(""),
		Core: meta.Core{
			VersionConstraint: "~1.0.0",
		},
	}

	err := r.Add(p)
	a.Nil(err)

	dep := meta.Dependency{
		ID:         "nori/test",
		Constraint: ">=1.0.0",
	}

	found := r.Resolve(dep)

	a.NotNil(found)
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
