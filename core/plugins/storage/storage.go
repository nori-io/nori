package storage

import (
	"sync"

	"errors"

	"strings"

	go_config "github.com/cheebo/go-config"
	"github.com/secure2work/nori/core/entities"
	"github.com/sirupsen/logrus"
)

type noriCoreStorage struct {
	Source string
	Log    *logrus.Logger
}

type NoriCoreStorage interface {
	GetInstallations() ([]entities.PluginMeta, error)
	SaveInstallation(meta entities.PluginMeta) error
	RemoveInstallation(id string) error
}

var instance NoriCoreStorage
var once sync.Once

func GetNoriCoreStorage(cfg go_config.Config, log *logrus.Logger) NoriCoreStorage {
	once.Do(func() {
		storageType := cfg.String("nori.storage.type")
		if len(storageType) == 0 {
			log.Error(errors.New("nori.storage.type not defined"))
			return
		}
		if strings.ToLower(storageType) == "none" {
			instance, _ = getNoneStorage()
			return
		}

		storageSource := cfg.String("nori.storage.source")
		if len(storageSource) == 0 {
			log.Error(errors.New("nori.storage.source not defined"))
		}

		var storage NoriCoreStorage
		var err error

		switch storageType {
		case "mysql":
			storage, err = getMySqlStorage(noriCoreStorage{
				Source: storageSource,
				Log:    log,
			})
			break
			//case "postgresql":
			//	// @todo implement
			//case "file":
			//	// @todo implement
		default:
			log.Error(errors.New("unknown nori.storage.type: " + storageType))
		}
		if err != nil {
			log.Error(err)
			return
		}
		instance = storage
	})
	return instance
}
