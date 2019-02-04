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
	"github.com/nori-io/nori-common/interfaces"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori-common/plugin"
	"github.com/nori-io/nori/core/errors"
	"github.com/sirupsen/logrus"
)

type registry struct {
	registryManager RegistryManager
	log             *logrus.Logger
	configManager   config.Manager
}

func NewRegistry(rm RegistryManager, cm config.Manager, logger *logrus.Logger) plugin.Registry {
	return registry{
		log:             logger,
		registryManager: rm,
		configManager:   cm,
	}
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

func (r registry) Auth() (interfaces.Auth, error) {
	item, err := r.registryManager.GetInterface(meta.Auth)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Auth)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.Auth,
		}
	}
	return i, nil
}

func (r registry) Authorize() (interfaces.Authorize, error) {
	item, err := r.registryManager.GetInterface(meta.Authorize)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Authorize)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.Authorize,
		}
	}
	return i, nil
}

func (r registry) Cache() (interfaces.Cache, error) {
	item, err := r.registryManager.GetInterface(meta.Cache)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Cache)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.Cache,
		}
	}
	return i, nil
}

func (r registry) Config() config.Manager {
	return r.configManager
}

func (r registry) Http() (interfaces.Http, error) {
	item, err := r.registryManager.GetInterface(meta.HTTP)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Http)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTP,
		}
	}
	return i, nil
}

func (r registry) HTTPTransport() (interfaces.HTTPTransport, error) {
	item, err := r.registryManager.GetInterface(meta.HTTPTransport)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.HTTPTransport)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTPTransport,
		}
	}
	return i, nil
}

func (r registry) Logger(meta meta.Meta) *logrus.Logger {
	return r.log.WithField("plugin", meta.Id().String()).Logger
}

func (r registry) Mail() (interfaces.Mail, error) {
	item, err := r.registryManager.GetInterface(meta.Mail)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Mail)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.Mail,
		}
	}
	return i, nil
}

func (r registry) PubSub() (interfaces.PubSub, error) {
	item, err := r.registryManager.GetInterface(meta.PubSub)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.PubSub)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.PubSub,
		}
	}
	return i, nil
}

func (r registry) Session() (interfaces.Session, error) {
	item, err := r.registryManager.GetInterface(meta.Session)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Session)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTPTransport,
		}
	}
	return i, nil
}

func (r registry) Sql() (interfaces.SQL, error) {
	item, err := r.registryManager.GetInterface(meta.SQL)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.SQL)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTPTransport,
		}
	}
	return i, nil
}

func (r registry) Templates() (interfaces.Templates, error) {
	item, err := r.registryManager.GetInterface(meta.Templates)
	if err != nil {
		return nil, err
	}
	i, ok := item.(interfaces.Templates)
	if !ok {
		return nil, errors.InterfaceAssertError{
			Interface: meta.HTTPTransport,
		}
	}
	return i, nil
}
