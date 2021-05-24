package nori

import (
	"context"

	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/common/v5/pkg/domain/plugin"
	"github.com/nori-io/common/v5/pkg/domain/registry"
	"github.com/nori-io/nori/pkg/nori/domain/enum"
)

type Nori interface {
	Add(p plugin.Plugin) error
	Remove(p plugin.Plugin) error

	Init(ctx context.Context, id meta.ID) error
	InitAll(ctx context.Context) error

	Start(ctx context.Context, id meta.ID) error
	StartAll(ctx context.Context) error

	Stop(ctx context.Context, id meta.ID) error
	StopAll(ctx context.Context) error

	Install(ctx context.Context, p plugin.Plugin) error
	UnInstall(ctx context.Context, id meta.ID) error

	GetByFilter(filter Filter) []meta.ID
	GetPluginVariables(id meta.ID) []registry.Variable

	GetState(id meta.ID) (enum.State, error)
}

type Filter struct {
	State enum.State
}
