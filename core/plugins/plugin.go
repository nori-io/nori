package plugins

import (
	"context"

	"github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/plugins/meta"
)

type Plugin interface {
	GetInstance() interface{}
	GetMeta() meta.Meta
	Init(ctx context.Context, config config.Manager) error
	Install(ctx context.Context, registry Registry) error
	Start(ctx context.Context, registry Registry) error
	Stop(ctx context.Context, registry Registry) error
	UnInstall(ctx context.Context, registry Registry) error
}
