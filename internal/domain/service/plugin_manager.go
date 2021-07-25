package service

import (
	"context"

	"github.com/nori-io/common/v5/pkg/domain/meta"
	"github.com/nori-io/nori/pkg/nori/domain/entity"
	"github.com/nori-io/nori/internal/domain/enum"
)

type PluginManager interface {
	Enable(ctx context.Context, id meta.ID) error
	Disable(ctx context.Context, id meta.ID) error

	Install(ctx context.Context, id meta.ID) error
	UnInstall(ctx context.Context, id meta.ID) error

	Start(ctx context.Context, id meta.ID) error
	Stop(ctx context.Context, id meta.ID) error

	StartAll(ctx context.Context) error
	StopAll(ctx context.Context) error

	GetByFilter(filter GetByFilterData) ([]*entity.Plugin, error)
}

type GetByFilterData struct {
	State enum.State
}
