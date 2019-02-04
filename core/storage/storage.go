// Copyright © 2018 Nori info@nori.io
//
// This program is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation, either version 3
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package storage

import (
	"sync"

	"github.com/nori-io/nori-common/meta"

	"strings"

	"fmt"

	go_config "github.com/cheebo/go-config"
	"github.com/sirupsen/logrus"
)

type theStorage struct {
	Source string
	Log    *logrus.Logger
}

type Storage interface {
	GetPluginMetas() ([]meta.Meta, error)
	SavePluginMeta(meta meta.Meta) error
	DeletePluginMeta(id meta.ID) error
}

const (
	storageTypeNone  = "none"
	storageTypeMysql = "mysql"

	cfgStorageType   = "nori.storage.type"
	cfgStorageSource = "nori.storage.source"
)

var instance Storage
var once sync.Once

func GetStorage(cfg go_config.Config, log *logrus.Logger) Storage {
	once.Do(func() {
		storageType := cfg.String(cfgStorageType)
		if len(storageType) == 0 {
			storageType = storageTypeNone
		}

		storageType = strings.ToLower(storageType)

		if storageType == storageTypeNone {
			instance, _ = getNoneStorage()
			return
		}

		storageSource := cfg.String(cfgStorageSource)
		if len(storageSource) == 0 {
			log.Error(fmt.Errorf("%s not defined", cfgStorageSource))
		}

		var storage Storage
		var err error

		switch storageType {
		case storageTypeMysql:
			storage, err = getMySqlStorage(storageSource, log)
			break
		default:
			log.Error(fmt.Errorf("unknown %s: %s", cfgStorageType, storageType))
		}
		if err != nil {
			log.Error(err)
			return
		}
		instance = storage
	})
	return instance
}
