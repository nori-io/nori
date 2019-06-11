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

package plugins

import (
	"github.com/nori-io/nori-common/config"
	"github.com/nori-io/nori-common/logger"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori-common/plugin"
	"github.com/nori-io/nori/core/errors"
)

type registry struct {
	registryManager RegistryManager
	log             logger.Logger
	configManager   config.Manager
}

func NewRegistry(rm RegistryManager, cm config.Manager, logger logger.Logger) plugin.Registry {
	return registry{
		log:             logger,
		registryManager: rm,
		configManager:   cm,
	}
}

func (r registry) Config() config.Manager {
	return r.configManager
}

func (r registry) Interface(i meta.Interface) (interface{}, error) {
	return r.registryManager.GetInterface(i)
}

func (r registry) Logger(meta meta.Meta) logger.Writer {
	//return r.log.WithFields(LogFieldsMeta(meta))
	return r.log
}

func (r registry) Resolve(dep meta.Dependency) (interface{}, error) {
	plugin := r.registryManager.Resolve(dep)
	if plugin != nil {
		return plugin.Instance(), nil
	}
	return nil, errors.DependencyNotFound{
		Dependency: dep,
	}
}
