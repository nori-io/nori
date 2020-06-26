package manager

import (
	"context"

	"github.com/nori-io/nori-common/v2/meta"
	"github.com/nori-io/nori/internal/domain/entity"
)

type Plugin interface {
	//
	Register(plugin *entity.Plugin) error
	//
	RegisterAll(plugins []*entity.Plugin) error
	//
	UnRegister(p *entity.Plugin) error

	//
	Install(ctx context.Context, plugin *entity.Plugin) error
	//
	UnInstall(ctx context.Context, plugin *entity.Plugin) error

	//
	GetAll() []*entity.Plugin
	//
	GetByID(id meta.ID) (*entity.Plugin, error)

	//
	GetHooks() []*entity.Plugin
	//
	GetInstallable() []*entity.Plugin
	//
	GetRunning() []*entity.Plugin

	//
	Start(ctx context.Context, id meta.ID) error
	//
	StartAll(ctx context.Context) error
	//
	Stop(ctx context.Context, id meta.ID) error
	//
	StopAll(ctx context.Context) error
}
