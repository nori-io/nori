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

	"github.com/nori-io/nori/core/plugins/dependency"
	"github.com/nori-io/nori/core/storage"
	"github.com/nori-io/nori/version"

	commonCfg "github.com/nori-io/nori-common/config"
	"github.com/nori-io/nori-common/meta"
	"github.com/nori-io/nori-common/plugin"
	"github.com/nori-io/nori/core/errors"

	"github.com/sirupsen/logrus"
)

type Manager interface {
	AddFile(path string) (plugin.Plugin, error)
	AddDir(paths []string) error

	Install(id meta.ID, ctx context.Context) error

	Meta(id meta.ID) (meta.Meta, error)
	Metas() []meta.Meta

	Start(id meta.ID, ctx context.Context) error
	StartAll(ctx context.Context) error

	Stop(id meta.ID, ctx context.Context) error
	StopAll(ctx context.Context) error

	UnInstall(id meta.ID, ctx context.Context) error
}

func NewManager(
	storage storage.Storage,
	configManager commonCfg.Manager,
	version version.Version,
	pluginExtractor PluginExtractor,
	logger *logrus.Logger,
) Manager {
	// @todo make as func param
	rm := NewRegistryManager(
		configManager,
		logger.WithField("component", "RegistryManager").Logger)

	return &manager{
		files:         map[string]meta.ID{},
		configManager: configManager,

		// @todo make as func param
		depManager:      dependency.NewManager(),
		pluginExtractor: pluginExtractor,
		regManager:      rm,

		// @todo make as func param
		registry: NewRegistry(rm, configManager, logger),
		storage:  storage,
		version:  version,
		log:      logger,
	}
}

type manager struct {
	files           FileTable
	configManager   commonCfg.Manager
	depManager      dependency.Manager
	pluginExtractor PluginExtractor
	regManager      RegistryManager
	registry        plugin.Registry
	storage         storage.Storage
	version         version.Version
	log             *logrus.Logger
}

func (m *manager) AddFile(path string) (plugin.Plugin, error) {
	p, err := m.pluginExtractor.Get(path)
	if err != nil {
		m.log.Error(err)
		return nil, err
	}

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
	if _, ok := p.(plugin.Installable); ok {
		// @todo check installed or not
		// if not installed - then
	}

	// add to dependency manager
	err = m.depManager.Add(p.Meta())
	if err != nil {
		return nil, err
	}

	err = m.regManager.Add(p)
	if err != nil {
		m.depManager.Remove(p.Meta().Id())
		return nil, err
	}

	m.files[path] = p.Meta().Id()

	return p, nil
}

func (m *manager) AddDir(paths []string) error {
	files, err := m.pluginExtractor.Files(paths)
	if err != nil {
		return err
	}

	for _, file := range files {
		// load plugin
		p, err := m.AddFile(file)
		if err != nil {
			m.log.Error(err)
			continue
		}

		m.log.Infof(
			"Found '%s' implements '%s' by '%s'",
			p.Meta().Id().String(),
			p.Meta().GetInterface(),
			p.Meta().GetAuthor().Name,
		)
	}
	return nil
}

func (m *manager) Install(id meta.ID, ctx context.Context) error {
	// @todo check depManager for dependencies
	p, err := m.regManager.Get(id)
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
	p, err := m.regManager.Get(id)
	if err != nil {
		return nil, err
	}
	return p.Meta(), nil
}

func (m *manager) Metas() []meta.Meta {
	var metas []meta.Meta
	for _, p := range m.regManager.Plugins() {
		metas = append(metas, p.Meta())
	}
	return metas
}

func (m *manager) Start(id meta.ID, ctx context.Context) error {
	// @todo check depManager for dependencies
	p, err := m.regManager.Get(id)
	if err != nil {
		return err
	}

	// @todo
	// - all deps must be resolved
	// - all deps must be already started
	// -- start non-started deps

	err = p.Init(ctx, m.configManager)
	if err != nil {
		return err
	}
	return p.Start(ctx, m.registry)
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
	// @todo stop dependent plugins
	p, err := m.regManager.Get(id)
	if err != nil {
		return err
	}

	return p.Stop(ctx, m.registry)
}

func (m *manager) StopAll(ctx context.Context) error {
	// todo stop running plugins in reverse order
	plugins := m.regManager.Plugins()
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
	p, err := m.regManager.Get(id)
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
	return nil
}
