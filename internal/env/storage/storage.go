package storage

import (
	"strings"

	"github.com/nori-io/common/v5/pkg/domain/storage"
	log "github.com/nori-io/logger"
	"github.com/nori-io/nori/internal/config"
	"github.com/nori-io/nori/internal/env/storage/bolt"
	"github.com/nori-io/nori/internal/env/storage/memory"
	errors2 "github.com/nori-io/nori/pkg/nori/domain/errors"
)

func New(cfg *config.Config) (storage.Storage, error) {
	var (
		store storage.Storage
		err   error
	)

	log.L().Info(cfg.Nori.Storage.DSN)
	// storage config
	if len(cfg.Nori.Storage.DSN) == 0 {
		return nil, errors2.ConfigParamUndefinedError{
			Param: "nori.storage.dsn",
		}
	}

	// schema://...
	parts := strings.Split(cfg.Nori.Storage.DSN, "://")
	switch parts[0] {
	case "mem":
		store, err = memory.New()
	case "bolt":
		if len(parts) < 2 {
			return nil, errors2.ConfigFormatError{
				Param: cfg.Nori.Storage.DSN,
			}
		}
		store, err = bolt.New(parts[1], 0666)
	default:
		return nil, errors2.ConfigFormatError{
			Param: cfg.Nori.Storage.DSN,
		}
	}

	if err != nil {
		return nil, err
	}

	return store, nil
}
