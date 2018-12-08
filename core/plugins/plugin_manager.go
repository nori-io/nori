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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"plugin"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/cheebo/go-config"
	cfgmanager "github.com/secure2work/nori/core/config"
	"github.com/secure2work/nori/core/entities"
	"github.com/secure2work/nori/core/plugins/interfaces"
	"github.com/secure2work/nori/core/storage"
)

type PluginManager interface {
	Load(dir []string) error
	LoadPlugin(filePath string) PluginEntry
	Plugins() map[string]PluginEntry
	Run(ctx context.Context, installed []entities.PluginMeta) error
	Install(id string) error
	UnInstall(id string) error
}

type manager struct {
	configManager interfaces.ConfigManager
	files         map[string]error
	log           *logrus.Logger
	plugins       map[string]PluginEntry
	registry      PluginRegistry
	storage       storage.NoriStorage
}

var instance *manager
var once sync.Once

func GetPluginManager(storage storage.NoriStorage, log *logrus.Logger, config go_config.Config) PluginManager {
	once.Do(func() {
		instance = &manager{
			configManager: cfgmanager.NewConfigManager(config),
			files:         map[string]error{},
			log:           log,
			plugins:       map[string]PluginEntry{},
			storage:       storage,
		}
		instance.registry = GetPluginRegistry(instance, log, instance.configManager)
	})
	return instance
}

func (m *manager) Plugins() map[string]PluginEntry {
	return m.plugins
}

func (m *manager) Load(sources []string) (err error) {
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
			entry := m.LoadPlugin(filePath)
			if entry == nil {
				continue
			}

			m.log.Infof(
				"Found: '%s:%s' by '%s'",
				entry.GetMeta().GetId(),
				entry.GetMeta().GetVersion(),
				entry.GetMeta().GetAuthor(),
			)
		}
	}
	return nil
}

func (m *manager) LoadPlugin(filePath string) PluginEntry {
	pluginFile, err := plugin.Open(filePath)
	if err != nil {
		m.log.WithField("file", filePath).Error(err)
		m.log.WithField("file", filePath).Error(PluginOpenError.Error())
		return nil
	}

	instance, err := pluginFile.Lookup("Plugin")
	if err != nil {
		m.log.WithField("file", filePath).Error(PluginLookupError.Error())
		return nil
	}

	plug, ok := instance.(Plugin)
	if !ok {
		m.log.WithField("file", filePath).Error(PluginInterfaceError.Error())
		return nil
	}

	var iface string
	switch plug.GetMeta().GetInterface() {
	case entities.Custom:
		iface = plug.GetMeta().GetId()
	default:
		iface = strings.ToLower(plug.GetMeta().GetInterface().String())
	}

	if len(iface) == 0 || iface == "unknown" {
		m.log.WithField("file", filePath).Error(PluginImplementedInterfaceError.Error())
		return nil
	}

	entry := NewPluginEntry(plug, filePath)

	m.plugins[iface] = entry

	return entry
}

func (m *manager) Install(id string) error {
	var plugin PluginEntry
	var ok bool
	if plugin, ok = m.plugins[id]; !ok {
		return &PluginNotFound{
			PluginId: id,
		}
	}
	// @todo check if already installed

	// @todo Check dependencies or not (?)

	// install plugin
	m.log.Infof("installing %s", id)
	err := plugin.Install(context.Background(), m.registry)
	if err != nil {
		m.log.Infof("can't install %s: ", id, err.Error())
		return err
	}

	// save plugin meta into core storage
	err = m.storage.SavePluginMeta(plugin.GetMeta())
	if err != nil {
		m.log.Infof("can't install %s: ", id, err.Error())
		return err
	}
	m.log.Infof("successfully installed %s", id)

	// @todo potential problem: different context on start here and in PluginManager.Run
	m.log.Infof("starting %s", id)
	err = plugin.Start(context.Background(), m.registry)
	if err != nil {
		m.log.Infof("error on start %s: %s", id, err.Error())
		return err
	}
	m.log.Infof("successfully started %s", id)

	return nil
}

func (m *manager) UnInstall(id string) error {
	var plugin PluginEntry
	var ok bool

	// @todo potential problem: different context on start here and in PluginManager.Run
	ctx := context.Background()

	if plugin, ok = m.plugins[id]; ok {
		// @todo add flag to stop and delete dependent plugins
		// Check dependencies
		var dependencies []string
		version := plugin.GetMeta().GetVersion()
		for _, p := range m.plugins {
			if ok, _ := p.isDependent(id, version); ok {
				dependencies = append(dependencies, fmt.Sprintf("%s:%s", p.GetMeta().GetId(), p.GetMeta().GetVersion()))
			}
		}
		if len(dependencies) > 0 {
			return &PluginHasDependentPlugins{
				PluginId:     id,
				Dependencies: dependencies,
			}
		}

		// Stop plugin
		err := plugin.Stop(ctx, m.registry)
		if err != nil {
			return err
		}

		// uninstall plugin
		err = plugin.UnInstall(ctx, m.registry)
		if err != nil {
			return err
		}
	}

	// delete plugin meta from core storage
	err := m.storage.DeletePluginMeta(id)
	if err != nil {
		return err
	}

	return nil
}

func (m *manager) Run(ctx context.Context, installed []entities.PluginMeta) error {
	var enabledPluginEntries []PluginEntry

	installedEntries := map[string]PluginEntry{}

	// only Custom PluginInterface must be installed, other Kinds do not need installation,
	// because they provide interfaces to service without any implementation
	for id, e := range m.Plugins() {
		// find plugin in installed and add to installedEntries
		meta := e.GetMeta()
		switch meta.GetInterface() {
		case entities.Custom:
			if isInstalled(meta, installed) {
				installedEntries[id] = e
			}
			break
		default:
			installedEntries[e.GetMeta().GetInterface().String()] = e
		}
	}

	errs := CheckDependencies(installedEntries)
	if len(errs) > 0 {
		for _, e := range errs {
			m.log.Error(e.Error())
		}
		return errors.New("unresolved dependencies")
	}

	enabledPluginEntries = SortPlugins(installedEntries)

	for _, pe := range enabledPluginEntries {
		err := pe.Init(ctx, m.configManager)
		if err != nil {
			m.log.WithFields(logrus.Fields{
				"p.id":   pe.GetMeta().GetId(),
				"p.name": pe.GetMeta().GetPluginName(),
				"call":   "plugin.Init",
			}).Error(err)
			return err
		}
		err = pe.Start(ctx, m.registry)
		if err != nil {
			m.log.WithFields(logrus.Fields{
				"p.id":   pe.GetMeta().GetId(),
				"p.name": pe.GetMeta().GetPluginName(),
				"call":   "plugin.Start",
			}).Error(err)
			return err
		}
		m.log.Infof("Started plugin %s", pe.GetMeta().GetId())
	}
	return nil
}
