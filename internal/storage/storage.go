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
	"github.com/nori-io/nori-common/logger"
	"github.com/nori-io/nori-common/meta"

	"strings"

	go_config "github.com/cheebo/go-config"
)

type Storage interface {
	Plugins() Plugins
}

type Plugins interface {
	All() ([]meta.Meta, error)
	Delete(meta.ID) error
	Get(meta.ID) (meta.Meta, error)
	IsInstalled(id meta.ID) (bool, error)
	Save(meta.Meta) error
}

const (
	storageTypeDummy = "dummy"
	storageTypeMysql = "mysql"

	configKeyStorageType   = "nori.storage.type"
	configKeyStorageSource = "nori.storage.source"
)

func NewStorage(cfg go_config.Config, log logger.Logger) (Storage, error) {
	storageType := cfg.String(configKeyStorageType)
	if len(storageType) == 0 {
		storageType = storageTypeDummy
	}

	storageSource := cfg.String(configKeyStorageSource)
	if len(storageSource) == 0 {
		return nil, UndefinedStorageSource{
			path: configKeyStorageSource,
		}
	}

	switch strings.ToLower(storageType) {
	case storageTypeDummy:
		return newDummyStorage()
	case storageTypeMysql:
		return newMySQLStorage(storageSource, log)
	default:
		return nil, UnknownStorageType{storageType: storageType}
	}
}
