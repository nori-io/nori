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
	"context"
	"fmt"

	"github.com/nori-io/nori-common/logger"

	"github.com/nori-io/nori/core/plugins/dependency"
	"github.com/nori-io/nori/core/plugins/types"
	"github.com/nori-io/nori/core/storage"
	"github.com/nori-io/nori/version"

	commonCfg "github.com/nori-io/nori-common/config"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori-common/plugin"
	"github.com/nori-io/nori/core/errors"
)

type Manager interface {
	AddFile(path string) (meta.Meta, error)
	AddDir(paths []string) error

	Install(id meta.ID, ctx context.Context) error

	Meta(id meta.ID) (meta.Meta, error)
	Metas(filter MetaFilter) types.MetaList

	Start(id meta.ID, ctx context.Context) error
	StartAll(ctx context.Context) error

	Stop(id meta.ID, ctx context.Context) error
	StopAll(ctx context.Context) error

	UnInstall(id meta.ID, ctx context.Context) error
}

type MetaFilter int

const (
	FilterRunnable MetaFilter = iota
	FilterInstallable
	FilterRunning
)

func NewManager(
	storage storage.Storage,
	configManager commonCfg.Manager,
	version version.Version,
	pluginExtractor PluginExtractor,
	logger logger.Logger,
) Manager {
	// @todo make as func param
	rm := NewRegistryManager(
		configManager,
		logger.WithField("component", "RegistryManager"))

	return &manager{
		files:         map[string]meta.Meta{},
		configManager: configManager,

		// @todo make as func param
		depManager:      dependency.NewManager(),
		pluginExtractor: pluginExtractor,
		registryManager: rm,

		// @todo make as func param
		registry: NewRegistry(rm, configManager, logger),
		storage:  storage,
		version:  version,
		log:      logger,
	}
}

type manager struct {
	files              FileTable
	installablePlugins types.PluginList
	runningPlugins     types.MetaList
	configManager      commonCfg.Manager
	depManager         dependency.Manager
	pluginExtractor    PluginExtractor
	registryManager    RegistryManager
	registry           plugin.Registry
	storage            storage.Storage
	version            version.Version
	log                logger.Logger
}

func (m *manager) AddFile(path string) (meta.Meta, error) {
	p, err := m.pluginExtractor.Get(path)
	if err != nil {
		return nil, err
	}

	// plugins files
	m.files[path] = p.Meta()

	// check needed Nori Core version
	cons, err := p.Meta().GetCore().GetConstraint()
	if err != nil {
		return nil, err
	}
	if !cons.Check(m.version.Version()) {
		return nil, errors.IncompatibleCoreVersion{
			Id:                 p.Meta().Id(),
			NeededCoreVersion:  p.Meta().GetCore().VersionConstraint,
			CurrentCoreVersion: m.version.Original(),
		}
	}

	// check installed or not
	// if plugin not installed then added plugin to list of installable plugins
	// and exit function
	if _, ok := p.(plugin.Installable); ok {
		installed, err := m.storage.Plugins().IsInstalled(p.Meta().Id())
		if err != nil {
			return nil, err
		}
		if !installed {
			m.installablePlugins.Add(p)
			return p.Meta(), nil
		}
	}

	// add to dependency manager
	err = m.depManager.Add(p.Meta())
	if err != nil {
		return nil, err
	}

	err = m.registryManager.Add(p)
	if err != nil {
		m.depManager.Remove(p.Meta().Id())
		return nil, err
	}

	return p.Meta(), nil
}

func (m *manager) AddDir(paths []string) error {
	files, err := m.pluginExtractor.Files(paths)
	if err != nil {
		return err
	}

	for _, file := range files {
		// load plugin
		mt, err := m.AddFile(file)
		if err != nil {
			m.log.Error(err)
			continue
		}

		m.log.Infof(
			"Found '%s' implements '%s' by '%s'",
			mt.Id().ID,
			mt.GetInterface(),
			mt.GetAuthor().Name,
		)
	}
	return nil
}

func (m *manager) Install(id meta.ID, ctx context.Context) error {
	// @todo check depManager for dependencies
	p, err := m.installablePlugins.Find(id)
	if err != nil {
		return err
	}
	installable, ok := p.(plugin.Installable)
	if !ok {
		return errors.NonInstallablePlugin{
			Id:   p.Meta().Id(),
			Path: m.files.Find(p.Meta().Id()),
		}
	}
	return installable.Install(ctx, m.registry)
}

func (m *manager) Meta(id meta.ID) (meta.Meta, error) {
	p, err := m.registryManager.Get(id)
	if err != nil {
		return nil, err
	}
	return p.Meta(), nil
}

func (m *manager) Metas(filter MetaFilter) types.MetaList {
	var metas types.MetaList

	switch filter {
	case FilterRunning:
		for _, v := range m.runningPlugins {
			metas.Add(v)
		}
	case FilterInstallable:
		for _, v := range m.installablePlugins {
			metas.Add(v.Meta())
		}
	case FilterRunnable:
		for _, p := range m.registryManager.Plugins() {
			metas = append(metas, p.Meta())
		}
	}

	return metas
}

func (m *manager) Start(id meta.ID, ctx context.Context) error {
	p, err := m.registryManager.Get(id)
	if err != nil {
		return err
	}

	// all dependencies must be resolvable
	// all dependencies must be started, otherwise start dependency
	var depErrs errors.DependenciesNotFound
	for _, dep := range p.Meta().GetDependencies() {
		did, err := m.depManager.Resolve(dep)
		if err != nil {
			depErrs.Add(p.Meta().Id(), dep)
			continue
		}

		if depMeta, err := m.runningPlugins.Find(did); depMeta == nil {
			if err != nil {
				return err
			}
			err = m.Start(did, ctx)
			if err != nil {
				return err
			}
		}
	}

	if depErrs.HasErrors() {
		return depErrs
	}

	err = p.Init(ctx, m.configManager)
	if err != nil {
		return err
	}

	var startErr error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		startErr = p.Start(ctx, m.registry)
	}()

	if recovered != nil {
		//return fmt.Errorf("%v", recovered)
	}

	if startErr != nil {
		return startErr
	}

	m.runningPlugins.Add(p.Meta())

	m.log.WithFields(LogFieldsMeta(p.Meta())).Infof("Plugin successfully started")

	return nil
}

func (m *manager) StartAll(ctx context.Context) error {
	// start plugins in dependency correct order
	pl, err := m.depManager.Sort()
	if err != nil {
		return err
	}

	for _, id := range pl {
		if err := m.Start(id, ctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) Stop(id meta.ID, ctx context.Context) error {
	if running, _ := m.runningPlugins.Find(id); running == nil {
		return nil
	}

	p, err := m.registryManager.Get(id)
	if err != nil {
		return err
	}

	// stop dependent plugins before stop the plugin
	for _, dep := range m.depManager.GetDependent(id) {
		depPlugin, err := m.registryManager.Get(dep)
		if err != nil {
			return err
		}
		//@todo collect errors
		if err := m.Stop(depPlugin.Meta().Id(), ctx); err != nil {
			m.log.WithFields(LogFieldsMeta(depPlugin.Meta())).Error(err)
		}
	}

	var stopErr error
	var recovered interface{}
	func() {
		defer func() {
			recovered = recover()
		}()
		stopErr = p.Stop(ctx, m.registry)
	}()

	m.runningPlugins.Remove(id)

	if recovered != nil {
		return fmt.Errorf("%s", recovered)
	}

	if stopErr != nil {
		return stopErr
	}

	m.log.WithFields(logger.Fields{
		"plugin_id":      id.ID,
		"plugin_version": id.Version,
		"interface":      p.Meta().GetInterface(),
	}).Info("Plugin successfully stopped")

	return nil
}

func (m *manager) StopAll(ctx context.Context) error {
	// todo stop running plugins in reverse order
	plugins := m.registryManager.Plugins()
	for i := len(plugins) - 1; i >= 0; i-- {
		p := plugins[i]
		err := m.Stop(p.Meta().Id(), ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *manager) UnInstall(id meta.ID, ctx context.Context) error {
	// @todo check depManager for dependent plugins
	p, err := m.registryManager.Get(id)
	if err != nil {
		return err
	}

	// stop plugin before uninstall it
	err = p.Stop(ctx, m.registry)
	if err != nil {
		return err
	}

	installable, ok := p.(plugin.Installable)
	if !ok {
		return errors.NonInstallablePlugin{
			Id:   p.Meta().Id(),
			Path: m.files.Find(p.Meta().Id()),
		}
	}
	err = installable.UnInstall(ctx, m.registry)
	if err != nil {
		m.log.Error(err)
		return err
	}
	err = m.storage.Plugins().Delete(id)
	if err != nil {
		m.log.Error(err)
		return err
	}

	m.installablePlugins.Add(p)

	return nil
}
