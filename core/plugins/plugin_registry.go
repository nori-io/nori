// Copyright © 2018 Secure2Work info@secure2work.com
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
	"github.com/sirupsen/logrus"

	"strings"

	"github.com/secure2work/nori/core/entities"
	"github.com/secure2work/nori/core/plugins/interfaces"
)

type PluginRegistry interface {
	Get(ns string) interface{}

	Auth() interfaces.Auth
	Authorize() interfaces.Authorize
	Cache() interfaces.Cache
	Config() interfaces.ConfigManager
	Http() interfaces.Http
	Logger(meta entities.PluginMeta) *logrus.Logger
	Mail() interfaces.Mail
	PubSub() interfaces.PubSub
	Session() interfaces.Session
	Sql() interfaces.SQL
	Templates() interfaces.Templates
}

func GetPluginRegistry(
	manager PluginManager,
	log *logrus.Logger,
	configManager interfaces.ConfigManager,
) PluginRegistry {
	return &registry{
		configManager: configManager,
		log:           log,
		manager:       manager,
	}
}

type registry struct {
	configManager interfaces.ConfigManager
	log           *logrus.Logger
	manager       PluginManager
}

func (r *registry) Get(ns string) interface{} {
	for n, p := range r.manager.Plugins() {
		if strings.ToLower(n) == strings.ToLower(ns) {
			return p.GetInstance()
		}
	}
	return nil
}

func (r *registry) Auth() interfaces.Auth {
	item := r.Get(entities.Auth.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Auth)
	if !ok {
		r.log.Error("Can't cast to Auth interface")
	}
	return i
}

func (r *registry) Authorize() interfaces.Authorize {
	item := r.Get(entities.Authorize.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Authorize)
	if !ok {
		r.log.Error("Can't cast to Authorize interface")
	}
	return i
}

func (r *registry) Cache() interfaces.Cache {
	item := r.Get(entities.Cache.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Cache)
	if !ok {
		r.log.Error("Can't cast to Cache interface")
	}
	return i
}

func (r *registry) Config() interfaces.ConfigManager {
	return r.configManager
}

func (r *registry) Http() interfaces.Http {
	item := r.Get(entities.HTTP.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Http)
	if !ok {
		r.log.Error("Can't cast to HTTP interface")
	}
	return i
}

func (r *registry) Logger(meta entities.PluginMeta) *logrus.Logger {
	return r.log.WithFields(logrus.Fields{
		"p.id":   meta.GetId(),
		"p.name": meta.GetPluginName(),
	}).Logger
}

func (r *registry) Mail() interfaces.Mail {
	item := r.Get(entities.Mail.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Mail)
	if !ok {
		r.log.Error("Can't cast to Mail interface")
	}
	return i
}

func (r *registry) PubSub() interfaces.PubSub {
	item := r.Get(entities.PubSub.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.PubSub)
	if !ok {
		r.log.Error("Can't cast to PubSub interface")
	}
	return i
}

func (r *registry) Session() interfaces.Session {
	item := r.Get(entities.Session.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Session)
	if !ok {
		r.log.Error("Can't cast to Session interface")
	}
	return i
}

func (r *registry) Sql() interfaces.SQL {
	item := r.Get(entities.SQL.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.SQL)
	if !ok {
		r.log.Error("Can't cast to SQL interface")
	}
	return i
}

func (r *registry) Templates() interfaces.Templates {
	item := r.Get(entities.Templates.String())
	if item == nil {
		return nil
	}
	i, ok := item.(interfaces.Templates)
	if !ok {
		r.log.Error("Can't cast to Templates interface")
	}
	return i
}
