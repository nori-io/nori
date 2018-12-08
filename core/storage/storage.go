package storage

import (
	"sync"

	"strings"

	"fmt"

	go_config "github.com/cheebo/go-config"
	"github.com/secure2work/nori/core/entities"
	"github.com/sirupsen/logrus"
)

type noriStorage struct {
	Source string
	Log    *logrus.Logger
}

type NoriStorage interface {
	GetPluginMetas() ([]entities.PluginMeta, error)
	SavePluginMeta(meta entities.PluginMeta) error
	DeletePluginMeta(id string) error
}

const (
	storageTypeNone  = "none"
	storageTypeMysql = "mysql"

	cfgNoriStorageType   = "nori.storage.type"
	cfgNoriStorageSource = "nori.storage.source"
)

var instance NoriStorage
var once sync.Once

func GetNoriStorage(cfg go_config.Config, log *logrus.Logger) NoriStorage {
	once.Do(func() {
		storageType := cfg.String(cfgNoriStorageType)
		if len(storageType) == 0 {
			storageType = storageTypeNone
		}

		storageType = strings.ToLower(storageType)

		if storageType == storageTypeNone {
			instance, _ = getNoneStorage()
			return
		}

		storageSource := cfg.String(cfgNoriStorageSource)
		if len(storageSource) == 0 {
			log.Error(fmt.Errorf("%s not defined", cfgNoriStorageSource))
		}

		var storage NoriStorage
		var err error

		switch storageType {
		case storageTypeMysql:
			storage, err = getMySqlStorage(noriStorage{
				Source: storageSource,
				Log:    log,
			})
			break
		default:
			log.Error(fmt.Errorf("unknown %s: %s", cfgNoriStorageType, storageType))
		}
		if err != nil {
			log.Error(err)
			return
		}
		instance = storage
	})
	return instance
}
