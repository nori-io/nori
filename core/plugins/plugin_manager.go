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
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/secure2work/nori/core/config"
	cfgmanager "github.com/secure2work/nori/core/config/manager"
	"github.com/secure2work/nori/core/entities"
	"github.com/secure2work/nori/core/interfaces"
	"github.com/secure2work/nori/core/plugins/storage"
)

type PluginManager interface {
	Load(dir []string) error
	LoadPlugin(filePath string) PluginEntry
	Plugins() map[string]PluginEntry
	Run(ctx context.Context, registry PluginRegistry, installed []entities.PluginMeta) error
	Install(id string) error
	UnInstall(id string) error
}

type manager struct {
	storage    storage.NoriStorage
	cfgManager interfaces.ConfigManager
	log        *logrus.Logger
	plugins    map[string]PluginEntry
	files      map[string]error
}

var instance *manager
var once sync.Once

func GetPluginManager(storage storage.NoriStorage) PluginManager {
	once.Do(func() {
		instance = &manager{
			storage:    storage,
			cfgManager: cfgmanager.NewConfigManager(config.Config),
			log:        config.Log,
			plugins:    map[string]PluginEntry{},
			files:      map[string]error{},
		}
	})
	return instance
}

func (r *manager) Plugins() map[string]PluginEntry {
	return r.plugins
}

func (r *manager) Load(sources []string) (err error) {
	for _, dir := range sources {
		var dirs []os.FileInfo
		if dirs, err = ioutil.ReadDir(dir); err != nil {
			return err
		}
		for _, d := range dirs {
			if d.IsDir() {
				continue
			}
			if path.Ext(d.Name()) != ".so" {
				continue
			}

			// load plugin
			filePath := filepath.Join(dir, d.Name())
			entry := r.LoadPlugin(filePath)
			if entry == nil {
				continue
			}

			r.log.Infof(
				"Found: '%s:%s' by '%s'",
				entry.Plugin().GetMeta().GetId(),
				entry.Plugin().GetMeta().GetVersion(),
				entry.Plugin().GetMeta().GetAuthor(),
			)
		}
	}
	return nil
}

func (r *manager) LoadPlugin(filePath string) PluginEntry {
	pluginFile, err := plugin.Open(filePath)
	if err != nil {
		r.log.WithField("file", filePath).Error(err)
		r.log.WithField("file", filePath).Error(PluginOpenError.Error())
		return nil
	}

	instance, err := pluginFile.Lookup("Plugin")
	if err != nil {
		r.log.WithField("file", filePath).Error(PluginLookupError.Error())
		return nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		r.log.WithField("file", filePath).WithField("os", "macos").Error(PluginOpenError.Error())
		return nil
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		r.log.WithField("file", filePath).Error(PluginHashError.Error())
		return nil
	}
	f.Close()

	plug, ok := instance.(Plugin)
	if !ok {
		r.log.WithField("file", filePath).Error(PluginInterfaceError.Error())
		return nil
	}

	var name string
	pluginInterface := plug.GetMeta().GetInterface()
	switch pluginInterface {
	case entities.Custom:
		name = plug.GetMeta().GetId()
	default:
		name = strings.ToLower(pluginInterface.String())
	}

	if len(name) == 0 {
		r.log.WithField("file", filePath).Error(PluginNamespaceError.Error())
		return nil
	}

	entry := &pluginEntry{
		plugin:   instance,
		filePath: filePath,
		hash:     hasher.Sum(nil),
		weight:   -1,
	}

	r.plugins[name] = entry

	return entry
}

func (r *manager) Install(id string) error {
	var plugin PluginEntry
	var ok bool
	if plugin, ok = r.plugins[id]; !ok {
		return &PluginNotFound{
			PluginId: id,
		}
	}
	// @todo check if already installed

	// @todo Check dependencies or not (?)

	// install plugin
	r.log.Infof("installing %s", id)
	err := plugin.Plugin().Install(context.Background(), r)
	if err != nil {
		r.log.Infof("can't install %s: ", id, err.Error())
		return err
	}

	// save plugin meta into core storage
	err = r.storage.SaveInstallation(plugin.Plugin().GetMeta())
	if err != nil {
		r.log.Infof("can't install %s: ", id, err.Error())
		return err
	}
	r.log.Infof("successfully installed %s", id)

	// @todo potential problem: different context on start here and in PluginManager.Run
	r.log.Infof("starting %s", id)
	err = plugin.Plugin().Start(context.Background(), r)
	if err != nil {
		r.log.Infof("error on start %s: %s", id, err.Error())
		return err
	}
	r.log.Infof("successfully started %s", id)

	return nil
}

func (r *manager) UnInstall(id string) error {
	var plugin PluginEntry
	var ok bool

	// @todo potential problem: different context on start here and in PluginManager.Run
	ctx := context.Background()

	if plugin, ok = r.plugins[id]; ok {
		// @todo add flag to stop and delete dependent plugins
		// Check dependencies
		var dependencies []string
		version := plugin.Plugin().GetMeta().GetVersion()
		for _, p := range r.plugins {
			if ok, _ := p.isDependent(id, version); ok {
				dependencies = append(dependencies, fmt.Sprintf("%s:%s", p.Plugin().GetMeta().GetId(), p.Plugin().GetMeta().GetVersion()))
			}
		}
		if len(dependencies) > 0 {
			return &PluginHasDependentPlugins{
				PluginId:     id,
				Dependencies: dependencies,
			}
		}

		// Stop plugin
		err := plugin.Plugin().Stop(ctx, r)
		if err != nil {
			return err
		}

		// uninstall plugin
		err = plugin.Plugin().UnInstall(ctx, r)
		if err != nil {
			return err
		}
	}

	// save plugin meta into core storage
	err := r.storage.RemoveInstallation(id)
	if err != nil {
		return err
	}

	return nil
}

func (r *manager) Run(ctx context.Context, registry PluginRegistry, installed []entities.PluginMeta) error {
	var enabledPluginEntries []PluginEntry

	installedEntries := map[string]PluginEntry{}

	// only Custom PluginInterface must be installed, other Kinds do not need installation,
	// because they provide interfaces to service without any implementation
	for id, e := range r.Plugins() {
		// find plugin in installed and add to installedEntries
		meta := e.Plugin().GetMeta()
		switch meta.GetInterface() {
		case entities.Custom:
			if isInstalled(meta, installed) {
				installedEntries[id] = e
			}
			break
		default:
			installedEntries[e.Plugin().GetMeta().GetInterface().String()] = e
		}
	}

	errs := CheckDependencies(installedEntries)
	if len(errs) > 0 {
		for _, e := range errs {
			r.log.Error(e.Error())
		}
		return errors.New("unresolved dependencies")
	}

	enabledPluginEntries = SortPlugins(installedEntries)

	for _, pe := range enabledPluginEntries {
		err := pe.Plugin().Init(ctx, r.cfgManager)
		if err != nil {
			r.log.WithFields(logrus.Fields{
				"p.id":   pe.Plugin().GetMeta().GetId(),
				"p.name": pe.Plugin().GetMeta().GetPluginName(),
				"call":   "plugin.Init",
			}).Error(err)
			return err
		}
		err = pe.Plugin().Start(ctx, registry)
		if err != nil {
			r.log.WithFields(logrus.Fields{
				"p.id":   pe.Plugin().GetMeta().GetId(),
				"p.name": pe.Plugin().GetMeta().GetPluginName(),
				"call":   "plugin.Start",
			}).Error(err)
			return err
		}
		r.log.Infof("Started plugin %s", pe.Plugin().GetMeta().GetId())
	}
	return nil
}

func (r *manager) Get(ns string) interface{} {
	for n, p := range r.plugins {
		if strings.ToLower(n) == strings.ToLower(ns) {
			return p.Plugin().GetInstance()
		}
	}
	return nil
}

func (r *manager) Auth() interfaces.Auth {
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

func (r *manager) Authorize() interfaces.Authorize {
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

func (r *manager) Cache() interfaces.Cache {
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

func (r *manager) Config() interfaces.ConfigManager {
	return r.cfgManager
}

func (r *manager) Http() interfaces.Http {
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

func (r *manager) Logger(meta entities.PluginMeta) *logrus.Logger {
	return r.log.WithFields(logrus.Fields{
		"p.id":   meta.GetId(),
		"p.name": meta.GetPluginName(),
	}).Logger
}

func (r *manager) Mail() interfaces.Mail {
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

func (r *manager) PubSub() interfaces.PubSub {
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

func (r *manager) Session() interfaces.Session {
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

func (r *manager) Sql() interfaces.SQL {
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

func (r *manager) Templates() interfaces.Templates {
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
