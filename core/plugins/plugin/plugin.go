package plugin

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
