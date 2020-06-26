package storage

import (
	"context"
	"strings"

	"github.com/nori-io/nori-common/v2/config"
	"github.com/nori-io/nori-common/v2/errors"
	"github.com/nori-io/nori-common/v2/logger"
	"github.com/nori-io/nori-common/v2/plugin"
	"github.com/nori-io/nori-common/v2/storage"
	"github.com/nori-io/nori/internal/domain/repository"
	"github.com/nori-io/nori/internal/env"
	"github.com/nori-io/nori/internal/service/storage/bolt"
	"github.com/nori-io/nori/internal/service/storage/memory"
	errors2 "github.com/nori-io/nori/pkg/errors"
	"go.uber.org/fx"
)

type Params struct {
	fx.In

	Env             *env.Env
	FileRepository  repository.FileRepository
	RegistryService plugin.Registry
	ConfigManager   config.Manager
	Logger          logger.Logger
}

func New(params Params) (storage.Storage, error) {
	var (
		s   storage.Storage
		err error
	)

	// default storage
	if len(params.Env.Config.Storage.DSN) == 0 {
		return memStorage()
	}

	// schema://...
	parts := strings.Split(params.Env.Config.Storage.DSN, "://")
	if len(parts) < 2 {
		return nil, errors2.ConfigFormatError{
			Param: params.Env.Config.Storage.DSN,
		}
	}
	switch parts[0] {
	case "bolt":
		s, err = boltStorage(parts[1])
	case "file":
		s, err = pluginStorage(parts[1], params.FileRepository, params.ConfigManager, params.RegistryService, params.Logger)
	default:
		return nil, errors2.ConfigFormatError{
			Param: params.Env.Config.Storage.DSN,
		}
	}

	if err != nil {
		return nil, err
	}

	return &Service{
		FileRepository: params.FileRepository,
		Env:            params.Env,
		Storage:        s,
	}, nil
}

func memStorage() (storage.Storage, error) {
	return memory.New()
}

func boltStorage(file string) (storage.Storage, error) {
	return bolt.New(file, 0666)
}

func pluginStorage(file string, fileRepo repository.FileRepository, cm config.Manager, r plugin.Registry, l logger.Logger) (storage.Storage, error) {
	f, err := fileRepo.File(file)
	if err != nil {
		return nil, err
	}
	p, err := fileRepo.Get(*f)
	if err != nil {
		return nil, err
	}
	// init
	ctx := context.Background()
	if err = p.Init(ctx, cm.Register(p.Meta().Id()), l); err != nil {
		return nil, err
	}
	// start
	if err = p.Start(ctx, r); err != nil {
		return nil, err
	}
	// instance
	instance, ok := p.Instance().(storage.Storage)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: storage.StorageInterface,
		}
	}
	return instance, nil
}
