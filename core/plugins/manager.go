package plugins

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	stdplugin "plugin"

	"github.com/secure2work/nori/core/plugins/errors"

	"github.com/secure2work/nori/version"

	"github.com/secure2work/nori/core/config"

	"github.com/secure2work/nori/core/storage"

	"github.com/secure2work/nori/core/plugins/meta"
	"github.com/secure2work/nori/core/plugins/plugin"

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
	storage storage.NoriStorage,
	registry RegistryManager,
	cfgManager config.Manager,
	version version.Version,
	logger *logrus.Logger,
) Manager {
	return &manager{
		files:      map[string]meta.ID{},
		plugins:    map[meta.ID]plugin.Plugin{},
		registry:   registry,
		log:        logger,
		cfgManager: cfgManager,
		storage:    storage,
		version:    version,
	}
}

type manager struct {
	files      FileTable
	log        *logrus.Logger
	plugins    map[meta.ID]plugin.Plugin
	registry   RegistryManager
	cfgManager config.Manager
	storage    storage.NoriStorage
	version    version.Version
}

func (m *manager) AddFile(path string) (plugin.Plugin, error) {
	file, err := stdplugin.Open(path)
	if err != nil {
		e := errors.FileOpenError{
			Path: path,
			Err:  err,
		}
		// @todo add err to error collector
		m.log.WithField("file", path).Error(e.Error())
		return nil, e
	}

	instance, err := file.Lookup("Plugin")
	if err != nil {
		e := errors.LookupError{
			Path: path,
			Err:  err,
		}
		// @todo add err to error collector
		m.log.WithField("file", path).Error(e.Error())
		return nil, e
	}

	p, ok := instance.(plugin.Plugin)
	if !ok {
		e := errors.TypeAssertError{
			Path: path,
		}
		// @todo add err to error collector
		m.log.WithField("file", path).Error(e.Error())
		return nil, e
	}

	// check needed Nori version
	cons, err := p.Meta().GetCore().GetConstraint()
	if err != nil {
		// @todo add err to error collector
		return nil, err
	}
	if !cons.Check(m.version.Version()) {
		// @todo add err to error collector
		return nil, errors.IncompatibleCoreVersion{
			Id:                 p.Meta().Id(),
			NeededCoreVersion:  p.Meta().GetCore().VersionConstraint,
			CurrentCoreVersion: m.version.Original(),
		}
	}

	err = m.registry.Add(p)
	if err != nil {
		// @todo add err to error collector
		return nil, err
	}

	m.plugins[p.Meta().Id()] = p
	m.files[path] = p.Meta().Id()

	return p, nil
}

func (m *manager) AddDir(paths []string) error {
	var err error
	for _, dir := range paths {
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
			p, err := m.AddFile(filePath)
			if err != nil {
				m.log.Error(err)
				continue
			}

			m.log.Infof(
				"Found '%s' by '%s'",
				p.Meta().Id().String(),
				p.Meta().GetAuthor().Name,
			)
		}
	}
	return nil
}

func (m *manager) Install(id meta.ID, ctx context.Context) error {
	p, ok := m.plugins[id]
	if !ok {
		return errors.NotFound{ID: id}
	}
	installable, ok := p.(plugin.Installable)
	if !ok {
		return errors.NonInstallablePlugin{
			Id:   p.Meta().Id(),
			Path: m.files.Find(p.Meta().Id()),
		}
	}
	return installable.Install(ctx, m.registry.Registry())
}

func (m *manager) Meta(id meta.ID) (meta.Meta, error) {
	p, ok := m.plugins[id]
	if !ok {
		return nil, errors.NotFound{
			ID: id,
		}
	}
	return p.Meta(), nil
}

func (m *manager) Metas() []meta.Meta {
	var meta []meta.Meta
	for _, p := range m.plugins {
		meta = append(meta, p.Meta())
	}
	return meta
}

func (m *manager) Start(id meta.ID, ctx context.Context) error {
	p, ok := m.plugins[id]
	if !ok {
		return errors.NotFound{ID: id}
	}
	// check:
	// - all deps must be resolved
	// - all deps must be already started
	// -- start non-started deps
	err := p.Init(ctx, m.cfgManager)
	if err != nil {
		return err
	}
	return p.Start(ctx, nil)
}

func (m *manager) StartAll(ctx context.Context) error {
	// todo start plugins in topological order
	pl, err := m.registry.OrderedPluginList()
	if err != nil {
		return err
	}

	for _, p := range pl {
		err := p.Init(ctx, m.cfgManager)
		if err != nil {
			return err
		}
		err = p.Start(ctx, m.registry.Registry())
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *manager) Stop(id meta.ID, ctx context.Context) error {
	p, ok := m.plugins[id]
	if !ok {
		return errors.NotFound{ID: id}
	}
	return p.Stop(ctx, nil)
}

func (m *manager) StopAll(ctx context.Context) error {
	// todo stop plugins in topological order
	return nil
}

func (m *manager) UnInstall(id meta.ID, ctx context.Context) error {
	p, ok := m.plugins[id]
	if !ok {
		return errors.NotFound{ID: id}
	}
	installable, ok := p.(plugin.Installable)
	if !ok {
		return errors.NonInstallablePlugin{
			Id:   p.Meta().Id(),
			Path: m.files.Find(p.Meta().Id()),
		}
	}
	err := installable.UnInstall(ctx, m.registry.Registry())
	if err != nil {
		m.log.Error(err)
		return err
	}
	err = m.storage.DeletePluginMeta(id)
	if err != nil {
		m.log.Error(err)
		return err
	}
	return nil
}
