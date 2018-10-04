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

type NoriStorage interface {
	GetInstallations() ([]entities.PluginMeta, error)
	SaveInstallation(meta entities.PluginMeta) error
	RemoveInstallation(id string) error
}

const (
	storageTypeNone  = "none"
	storageTypeMysql = "mysql"
	storageTypeFile  = "file"
)

var instance NoriStorage
var once sync.Once

func GetNoriStorage(cfg go_config.Config, log *logrus.Logger) NoriStorage {
	once.Do(func() {
		storageType := cfg.String("nori.storage.type")
		if len(storageType) == 0 {
			storageType = storageTypeNone
		}

		storageType = strings.ToLower(storageType)

		// if storageType == none then return NoneStorage
		if storageType == storageTypeNone {
			instance, _ = getNoneStorage()
			return
		}

		storageSource := cfg.String("nori.storage.source")
		if len(storageSource) == 0 {
			log.Error(errors.New("nori.storage.source not defined"))
		}

		var storage NoriStorage
		var err error

		switch storageType {
		case storageTypeMysql:
			storage, err = getMySqlStorage(noriCoreStorage{
				Source: storageSource,
				Log:    log,
			})
			break
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
